package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"Loop_backend/internal/dto"
	"Loop_backend/internal/response"
	"Loop_backend/internal/services"
)

type AuthHandler struct {
	userService services.UserService
	authService services.AuthService
}

func NewAuthHandler(userService services.UserService, authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Create a user
	user, err := h.userService.CreateUser(req.Email, req.Username)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Hash and add the password to the password table
	if err := h.authService.RegisterUserPassword(user.ID, req.Password); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Create session for the new user
	session, err := h.authService.CreateSession(user.ID)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := dto.AuthResponse{
		UserID:       session.UserID,
		AccessToken:  session.Token,
		ExpiresAt:    session.ExpiresAt.Format(time.RFC3339),
	}

	response.RespondWithJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user_id, err := h.authService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
        fmt.Println(err)
		response.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	session, err := h.authService.CreateSession(user_id)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := dto.AuthResponse{
		UserID:       session.UserID,
		AccessToken:  session.Token,
		ExpiresAt:    session.ExpiresAt.Format(time.RFC3339),
	}

	response.RespondWithJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		response.RespondWithError(w, http.StatusUnauthorized, "No token provided")
		return
	}

	claims, err := h.authService.ValidateToken(authHeader)

	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, claims)
}

func (h *AuthHandler) RegisterRoutes(r *RouteRegister) {
	r.RegisterPublicRoute("/api/auth/register", h.Register)
	r.RegisterPublicRoute("/api/auth/login", h.Login)
	r.RegisterProtectedRoute("/api/auth/verify", h.Verify)
}
