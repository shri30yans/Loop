package dto

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type AuthResponse struct {
    UserID       string `json:"user_id"`
    AccessToken  string `json:"access_token"`
    ExpiresAt    string `json:"expires_at"`
}

type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Username string `json:"username" validate:"required,min=3"`
    Password string `json:"password" validate:"required,min=6"`
}
