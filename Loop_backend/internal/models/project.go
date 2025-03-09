package models

import (
    "errors"
    "time"
)

// Project represents a project in the system
type Project struct {
    ID           string       `json:"id"`
    OwnerID      string       `json:"owner_id"`
    Title        string    `json:"title"`
    Description  string    `json:"description"`
    Introduction string    `json:"introduction"`
    Tags         []string  `json:"tags"`
    Sections     []Section `json:"sections"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// Section represents a project section
type Section struct {
    Index   string  `json:"index"` // starting from 1
    Title   string `json:"title"`
    Content string `json:"content"`
}

var (
    ErrInvalidProjectID    = errors.New("invalid project ID")
    ErrInvalidOwnerID      = errors.New("invalid owner ID")
    ErrInvalidTitle        = errors.New("invalid project title")
    ErrInvalidDescription  = errors.New("invalid project description")
)

// NewProject creates a new project instance with validation
func NewProject(ownerID string, title, description, introduction string, tags []string) (*Project, error) {
    if ownerID == "" {
        return nil, ErrInvalidOwnerID
    }
    if title == "" {
        return nil, ErrInvalidTitle
    }
    if description == "" {
        return nil, ErrInvalidDescription
    }

    now := time.Now()
    return &Project{
        OwnerID:      ownerID,
        Title:        title,
        Description:  description,
        Introduction: introduction,
        Tags:         tags,
        Sections:     make([]Section, 0),
        CreatedAt:    now,
        UpdatedAt:    now,
    }, nil
}
