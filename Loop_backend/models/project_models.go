package models

import (
	"time"
)

// ----------------------------------------------------------------------------
// Project Structures
// ----------------------------------------------------------------------------

// Project represents a project in the database.
type Project struct {
	ProjectID    int              `json:"project_id"`
	OwnerID      int              `json:"owner_id"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Introduction string           `json:"introduction"`
	Status       *string          `json:"status"`
	CreatedAt    time.Time        `json:"created_at"`
	Sections     []ProjectSection `json:"sections"`
	Tags         []string         `json:"tags"`
	Owner        UserDetails      `json:"owner"`
}

type T struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Introduction string `json:"introduction"`
	OwnerId      string `json:"owner_id"`
	Tags         string `json:"tags"`
}

// Feedback represents feedback for a project.
type Feedback struct {
	FeedbackID int    `json:"feedback_id"`
	ProjectID  int    `json:"project_id"`
	UserID     int    `json:"user_id"`
	Feedback   string `json:"feedback"`
}

// ProjectSection represents an update for a project.
type ProjectSection struct {
	Title         string `json:"title"`
	SectionNumber int    `json:"section_number"`
	Body          string `json:"body"`
	ProjectID     int    `json:"project_id"`
}

type ProjectTag struct {
	ProjectID      int    `json:"project_id"`
	TagDescription string `json:"tag_description"`
}

type ProjectsResponse struct {
	ProjectID   int       `json:"project_id"`
	OwnerID     int       `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Tags        []string  `json:"tags"`
}
