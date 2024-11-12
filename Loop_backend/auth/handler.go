package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user, err := CreateUser(req.Name, strings.ToLower(req.Email), string(hashedPassword))
	if err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			http.Error(w, "Email already exists", http.StatusConflict)
		} else {
			http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(RegisterResponse{
		UserID:         fmt.Sprintf("%d", user.ID),
		Email:          user.Email,
		HashedPassword: string(hashedPassword),
	})
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := GetUserByEmail(strings.ToLower(req.Email))
	if err != nil {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check if a session already exists
	existingSession, err := GetSessionByUserID(user.ID)
	if err == nil {
		// Return existing session
		json.NewEncoder(w).Encode(AuthResponse{
			UserID:       fmt.Sprintf("%d", user.ID),
			RefreshToken: existingSession.RefreshToken,
			ExpiresAt:    existingSession.ExpiresAt.Format(time.RFC3339),
		})
		return
	}

	// Create new session if no existing session is found
	session, err := CreateSession(user.ID)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AuthResponse{
		UserID:       fmt.Sprintf("%d", user.ID),
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt.Format(time.RFC3339),
	})
}

func HandleVerify(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
		return
	}

	// Remove Bearer prefix
	if len(token) > 7 && strings.ToUpper(token[:7]) == "BEARER " {
		token = token[7:]
	}
	fmt.Println("Verifying token", token)

	session, err := GetSessionByRefreshToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	fmt.Println(session)
	json.NewEncoder(w).Encode(map[string]string{
		"session_id": fmt.Sprintf("%d", session.SessionID),
	})
}

func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := GetSessionByRefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if session.ExpiresAt.Before(time.Now()) {
		DeleteSession(session.SessionID)
		http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(session.UserID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("X-Refresh-Token")
	if refreshToken != "" {
		if session, err := GetSessionByRefreshToken(refreshToken); err == nil {
			DeleteSession(session.SessionID)
		}
	}
	w.WriteHeader(http.StatusOK)
}
