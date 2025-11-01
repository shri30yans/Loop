package dto

import (
	"Loop_backend/internal/models"
	"github.com/google/uuid"
	"time"
)

type CreateProjectRequest struct {
    Title        string           `json:"title" validate:"required,min=3,max=100"`
    Description  string           `json:"description" validate:"required,min=10,max=500"`
    Status       string           `json:"status" validate:"omitempty,oneof=draft published archived"`
    Introduction string           `json:"introduction" validate:"required,min=10,max=1000"`
    Tags         []string         `json:"tags" validate:"omitempty,dive,min=2,max=30"`
    Sections     []models.Section `json:"sections" validate:"omitempty,dive"`
}

// ToModel converts the DTO to a Project model
func (r *CreateProjectRequest) ToModel(userID uuid.UUID) *models.Project {
    now := time.Now()

    // Set default status if not provided
    status := r.Status
    if status == "" {
        status = string(models.StatusPending)
    }

    // Initialize empty slices if not provided
    tags := r.Tags
    if tags == nil {
        tags = make([]string, 0)
    }

    sections := r.Sections
    if sections == nil {
        sections = make([]models.Section, 0)
    }

    return &models.Project{
        ProjectInfo: models.ProjectInfo{
            ProjectID:    uuid.New(),
            Title:        r.Title,
            Description:  r.Description,
            Status:       models.Status(status),
            Introduction: r.Introduction,
            Tags:         tags,
            CreatedAt:    now,
            UpdatedAt:    now,
            OwnerID:      userID,
        },
        Sections: sections,
    }
}

type ProjectResponse struct {
    ProjectID    uuid.UUID        `json:"projectId"`
    Title        string          `json:"title"`
    Description  string          `json:"description"`
    Status       string          `json:"status"`
    Introduction string          `json:"introduction"`
    Tags         []string        `json:"tags"`
    Sections     []models.Section `json:"sections"`
    CreatedAt    time.Time       `json:"createdAt"`
    UpdatedAt    time.Time       `json:"updatedAt"`
}

// FromModel converts a Project model to a ProjectResponse DTO
func FromModel(project *models.Project) *ProjectResponse {
    return &ProjectResponse{
        ProjectID:    project.ProjectID,
        Title:        project.Title,
        Description:  project.Description,
        Status:       string(project.Status),
        Introduction: project.Introduction,
        Tags:         project.Tags,
        Sections:     project.Sections,
        CreatedAt:    project.CreatedAt,
        UpdatedAt:    project.UpdatedAt,
    }
}
