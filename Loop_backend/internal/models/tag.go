package models

import (
	"time"

	"github.com/google/uuid"
)

// Tag represents a tag entity in the system
type Tag struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Embedding   []float64 `json:"-"`
	Vector      []float64 `json:"vector,omitempty"`
	UsageCount  int       `json:"usage_count"`
	Confidence  float64   `json:"confidence,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SetConfidence sets the confidence value for the tag
func (t *Tag) SetConfidence(confidence float64) {
	t.Confidence = confidence
}

type TagRelationship struct {
	Tag1ID        uuid.UUID `json:"tag1_id"`
	Tag2ID        uuid.UUID `json:"tag2_id"`
	Strength      float64   `json:"strength"`
	CoOccurrences int       `json:"co_occurrences"`
	LastUpdated   time.Time `json:"last_updated"`
}

type TagSuggestion struct {
	SuggestedTagID uuid.UUID `json:"suggested_tag_id"`
	Reason         string    `json:"reason"`
	Confidence     float64   `json:"confidence"`
}
