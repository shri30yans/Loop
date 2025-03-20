package models

import (
    "errors"
    "time"
    "github.com/google/uuid"
)

var (
    ErrInvalidProjectID    = errors.New("invalid project ID")
    ErrInvalidOwnerID      = errors.New("invalid owner ID")
    ErrInvalidTitle        = errors.New("invalid project title")
    ErrInvalidDescription  = errors.New("invalid project description")
)

type Status string

const (
	StatusPending   Status = "draft"
	StatusInProgress Status = "published"
	StatusCompleted  Status = "completed"
	StatusArchived   Status = "archived"
)

// ProjectInfo represents project metadata without sections
type ProjectInfo struct {
    ProjectID    string    `json:"id"`
    OwnerID      string    `json:"owner_id"`
    Title        string    `json:"title"`
    Description  string    `json:"description"`
    Status       Status    `json:"status"`
    Introduction string    `json:"introduction"`
    Tags         []string  `json:"tags"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// Project represents a full project in the system including sections
type Project struct {
    ProjectInfo
    Sections []Section `json:"sections"`
}

// ToProjectInfo converts a Project to ProjectInfo
func (p *Project) ToProjectInfo() *ProjectInfo {
    return &ProjectInfo{
        ProjectID:    p.ProjectID,
        OwnerID:      p.OwnerID,
        Title:        p.Title,
        Description:  p.Description,
        Status:       p.Status,
        Introduction: p.Introduction,
        Tags:         p.Tags,
        CreatedAt:    p.CreatedAt,
        UpdatedAt:    p.UpdatedAt,
    }
}

// Section represents a project section
type Section struct {
    Title   string `json:"title"`
    Body string `json:"body"`
}


// NewProject creates a new project instance with validation
func NewProject(ownerID, title, description, status, introduction string, tags []string, sections []Section) (*Project, error) {
    now := time.Now()
    return &Project{
        ProjectInfo: ProjectInfo{
            ProjectID:    uuid.NewString(),
            OwnerID:     ownerID,
            Title:       title,
            Description: description,
            Status:      Status(status),
            Introduction: introduction,
            Tags:        tags,
            CreatedAt:   now,
            UpdatedAt:   now,
        },
        Sections: sections,
    }, nil
}

// NewSection creates a new section with validation
func NewSection(title, body string) (*Section, error) {
	return &Section{
		Title:   title,
		Body: body,
	}, nil
}
