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

<<<<<<< HEAD
// ProjectInfo represents project metadata without sections
type ProjectInfo struct {
    ProjectID    string    `json:"id"`
    OwnerID      string    `json:"owner_id"`
=======
// Project represents a project in the system
type Project struct {
    ProjectID   string `json:"id"`
    OwnerID      string  `json:"owner_id"`
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
    Title        string    `json:"title"`
    Description  string    `json:"description"`
    Status       Status    `json:"status"`
    Introduction string    `json:"introduction"`
    Tags         []string  `json:"tags"`
<<<<<<< HEAD
=======
    Sections     []Section `json:"sections"`
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

<<<<<<< HEAD
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

=======
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
// Section represents a project section
type Section struct {
    Title   string `json:"title"`
    Body string `json:"body"`
}


// NewProject creates a new project instance with validation
<<<<<<< HEAD
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
=======
func NewProject(ownerID, title, description,status, introduction string, tags []string, sections []Section) (*Project, error) {
    now := time.Now()
    return &Project{
        ProjectID : uuid.NewString(),
        OwnerID:      ownerID,
        Title:        title,
        Description:  description,
        Status:       Status(status),
        Introduction: introduction,
        Tags:         tags,
        Sections:     sections,
        CreatedAt:    now,
        UpdatedAt:    now,
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
    }, nil
}

// NewSection creates a new section with validation
func NewSection(title, body string) (*Section, error) {
	return &Section{
		Title:   title,
		Body: body,
	}, nil
}
