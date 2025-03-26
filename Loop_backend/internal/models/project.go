package models

import (
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	StatusPending    Status = "draft"
	StatusInProgress Status = "published"
	StatusCompleted  Status = "completed"
	StatusArchived   Status = "archived"
)

type ProjectInfo struct {
	ProjectID    uuid.UUID `json:"id"`
	OwnerID      uuid.UUID `json:"owner_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Status       Status    `json:"status"`
	Introduction string    `json:"introduction"`
	Tags         []string  `json:"tags"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Section represents a project section
type Section struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Project represents a full project in the system including sections
type Project struct {
	ProjectInfo
	Sections []Section `json:"sections"`
}

// NewProject creates a new project instance with validation
func NewProject(ownerID, title, description, status, introduction string, tags []string, sections []Section) (*Project, error) {

	parsedOwnerID, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Project{
		ProjectInfo: ProjectInfo{
			ProjectID:    uuid.New(),
			OwnerID:      parsedOwnerID,
			Title:        title,
			Description:  description,
			Status:       Status(status),
			Introduction: introduction,
			Tags:         tags,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		Sections: sections,
	}, nil
}

// NewSection creates a new section with validation
func NewSection(title, content string) (*Section, error) {
	return &Section{
		Title:   title,
		Content: content,
	}, nil
}
