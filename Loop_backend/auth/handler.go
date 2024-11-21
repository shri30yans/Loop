package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	. "Loop/models"

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
		if errors.Is(err, errors.New("email already exists")) {
			http.Error(w, "Email already exists", http.StatusConflict)
		} else {
			http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(RegisterResponse{
		UserID:         user.ID,
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
		fmt.Println("Invalid User")
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		fmt.Println("Invalid Password")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check for existing session
	session, err := GetSessionByUserID(user.ID)
	if err != nil || session.ExpiresAt.Before(time.Now()) {
		// Delete expired session if it exists
		if session != nil {
			DeleteSession(session.SessionID)
		}

		session, err = CreateSession(user.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
	}

	jwtToken, err := GenerateJWT(user.ID)
	if err != nil {
		fmt.Println("Failed to generate JWT:", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	authResponse := AuthResponse{
		UserID:       user.ID,
		AccessToken:  jwtToken,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt.Format(time.RFC3339),
	}

	fmt.Printf("AuthResponse: %+v\n", authResponse)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse)
}

func HandleVerify(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== Starting Token Verification ===")

	claims, err := ValidateAndProcessToken(w, r)
	if err != nil {
		fmt.Printf("❌ Token validation failed: %v\n", err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Extract user_id
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	userID := int(userIDFloat)

	fmt.Printf("✅ Verified token for user ID: %d\n", userID)

	// Optionally fetch session details
	session, err := GetSessionByUserID(userID)
	if err != nil {
		fmt.Printf("❌ Failed to fetch session for user ID %d: %v\n", userID, err)
		http.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	// Respond with session details
	fmt.Printf("✅ Returning session details for user ID: %d\n", userID)
	json.NewEncoder(w).Encode(map[string]string{
		"session_id": fmt.Sprintf("%d", session.SessionID),
	})
}

func HandleEditPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		fmt.Println("no token provided")
	}

	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	session, err := GetSessionByRefreshToken(tokenString)
	if err != nil {
		fmt.Println("Invalid session")
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}
	userID := session.UserID

	user, err := GetUserByID(userID)
	if err != nil {
		fmt.Println("Invalid User")
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.CurrentPassword)); err != nil {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	if err := UpdateUserPassword(userID, string(hashedPassword)); err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password updated successfully",
	})
}

// func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
// 	var req struct {
// 		RefreshToken string `json:"refresh_token"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	session, err := GetSessionByRefreshToken(req.RefreshToken)
// 	if err != nil {
// 		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
// 		return
// 	}

// 	if session.ExpiresAt.Before(time.Now()) {
// 		DeleteSession(session.SessionID)
// 		http.Error(w, "Refresh token expired", http.StatusUnauthorized)
// 		return
// 	}

// 	token, err := GenerateJWT(session.UserID)
// 	if err != nil {
// 		http.Error(w, "Error generating token", http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(map[string]string{
// 		"token": token,
// 	})
// }

func HandleGetUserInfo(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Unauthorized: No user ID provided", http.StatusUnauthorized)
		return
	}

	user, err := GetUserInfoById(userID)
	if err != nil {

		http.Error(w, "Error fetching user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
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
