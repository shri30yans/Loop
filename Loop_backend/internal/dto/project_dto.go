package dto

import (
	"Loop_backend/internal/models"
)

type CreateProjectRequest struct {
	Title        string `json:"title" validate:"required"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	Introduction string `json:"introduction"`
	// OwnerID      string           `json:"owner_id" validate:"required"`
	Tags     []string         `json:"tags"`
	Sections []models.Section `json:"sections"`
}
