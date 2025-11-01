package models

import (
	"github.com/google/uuid"
	"time"
)

// TagRelationship represents a relationship between two tags
type TagRelationship struct {
	ID            uuid.UUID `json:"id"`
	SourceID      uuid.UUID `json:"source_id"`
	TargetID      uuid.UUID `json:"target_id"`
	Tag2Name      string    `json:"tag2_name"`
	Type          string    `json:"type"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	Strength      float64   `json:"strength"`
	CoOccurrences int       `json:"co_occurrences"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewTagRelationship creates a new tag relationship with default values
func NewTagRelationship() *TagRelationship {
	return &TagRelationship{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// BeforeCreate sets default values before creation
func (tr *TagRelationship) BeforeCreate() {
	if tr.ID == uuid.Nil {
		tr.ID = uuid.New()
	}
	if tr.CreatedAt.IsZero() {
		tr.CreatedAt = time.Now()
	}
	if tr.UpdatedAt.IsZero() {
		tr.UpdatedAt = tr.CreatedAt
	}
	if tr.CoOccurrences == 0 {
		tr.CoOccurrences = 1
	}
}

// BeforeUpdate updates the update timestamp
func (tr *TagRelationship) BeforeUpdate() {
	tr.UpdatedAt = time.Now()
}
