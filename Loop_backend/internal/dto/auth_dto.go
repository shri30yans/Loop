package dto

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type AuthResponse struct {
    UserID       string    `json:"user_id"`
    AccessToken  string `json:"access_token"`
    ExpiresAt    string `json:"expires_at"`
}
