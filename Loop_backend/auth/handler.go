package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Wrong request body. Expected format: {\"email\": \"user@example.com\", \"password\": \"yourpassword\"}", http.StatusBadRequest)
		return
	}

	//if !auth.isValidEmail(req.Email) || !auth.isValidPassword(req.Password) {
	//	http.Error(w, "Invalid email or password", http.StatusBadRequest)
	//	return
	//}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user, err := CreateUser(req.Name, req.Email, string(hashedPassword))
	if err != nil {
		//if auth.isDuplicateEmail(err) {
		//	http.Error(w, "Email already exists", http.StatusConflict)
		//	return
		//}
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User:  user,
	})
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := GetUserByEmail(req.Email)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User:  user,
	})
}
