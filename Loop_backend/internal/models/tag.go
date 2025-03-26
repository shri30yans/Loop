package models

import (
	"github.com/google/uuid"
	"time"
)

// Tag represents a project tag
type Tag struct {
	ID          uuid.UUID `json:"id"`
	ProjectID   uuid.UUID `json:"project_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	UsageCount  int       `json:"usage_count"`
	Embedding   []float64 `json:"embedding,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (t *Tag) BeforeCreate() {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}

	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	if t.UsageCount == 0 {
		t.UsageCount = 1
	}
}
