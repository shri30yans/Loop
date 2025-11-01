package handlers

import (
	"fmt"
	"net/http"
	"time"

	"Loop_backend/internal/dto"
	"Loop_backend/internal/middleware"
	"Loop_backend/internal/response"
	"Loop_backend/internal/services/auth"
	"Loop_backend/internal/services/user"
)

type AuthHandler struct {
	userService user.UserService
	authService auth.AuthService
}

func NewAuthHandler(userService user.UserService, authService auth.AuthService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		authService: authService,
	}
}

func (h *AuthHandler) RegisterRoutes(r RouteRegister) {
    // POST endpoints with request body validation
    r.RegisterPublicRouteWithMethod("/api/auth/register", http.MethodPost, middleware.ValidateRequest(h.Register, &dto.RegisterRequest{}))
    r.RegisterPublicRouteWithMethod("/api/auth/login", http.MethodPost, middleware.ValidateRequest(h.Login, &dto.LoginRequest{}))
    
    // GET endpoints without request body validation
    r.RegisterProtectedRouteWithMethod("/api/auth/verify", http.MethodGet, h.Verify)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	req, ok := middleware.GetDTO[*dto.RegisterRequest](r)
	if !ok {
response.RespondWithErrorDetails(w, http.StatusBadRequest, "Invalid request payload", map[string]interface{}{
    "reason":          "Failed to parse or validate request body",
    "expected_fields": "email, username, password",
})
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
		UserID:      session.UserID,
		AccessToken: session.Token,
		ExpiresAt:   session.ExpiresAt.Format(time.RFC3339),
	}

	response.RespondWithJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	req, ok := middleware.GetDTO[*dto.LoginRequest](r)
	if !ok {
response.RespondWithErrorDetails(w, http.StatusBadRequest, "Invalid request payload", map[string]interface{}{
    "reason":          "Failed to parse or validate request body",
    "expected_fields": "email, password",
})
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
		UserID:      session.UserID,
		AccessToken: session.Token,
		ExpiresAt:   session.ExpiresAt.Format(time.RFC3339),
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
	fmt.Println(claims)

	if err != nil {
		response.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusOK, claims)
}
