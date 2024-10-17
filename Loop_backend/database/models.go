package database

import (
	"github.com/google/uuid"
	"time"
)

// User represents a user in the database.
type User struct {
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	Location  string    `json:"location"`
	Bio       string    `json:"bio"`
}

// Project represents a project in the database.
type Project struct {
	ProjectID    uuid.UUID        `json:"project_id"`
	OwnerID      int              `json:"owner_id"`
	Title        string           `json:"title"`
	Introduction string           `json:"introduction"`
	Sections     []ProjectSection `json:"sections"`
	Description  string           `json:"description"`
	Status       string           `json:"status"`
	CreatedAt    time.Time        `json:"created_at"`
	Tags         string           `json:"tags"`
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

// ProjectUpdate represents an update for a project.
type ProjectSection struct {
	Title        string `json:"title"`
	UpdateNumber int    `json:"update_number"`
	Body         string `json:"body"`
	ProjectID    int    `json:"project_id"`
}

// Event represents an event in the database.
type Event struct {
	EventID int    `json:"event_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Company string `json:"company"`
}

// UserEventParticipation represents the participation of a user in an event.
type UserEventParticipation struct {
	UserID  int `json:"user_id"`
	EventID int `json:"event_id"`
}
