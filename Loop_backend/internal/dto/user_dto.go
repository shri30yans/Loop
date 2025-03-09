package dto

type RegisterRequest struct {
    Email    string `json:"email"`
    Username string `json:"username"`
}