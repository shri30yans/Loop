package models

import (
	"time"
)

// Tag represents a project tag with its vector representation and metadata
type Tag struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Vector      []float64 `json:"-"` // Vector representation for similarity calculations
	UsageCount  int       `json:"usage_count"`
	Confidence  float64   `json:"confidence,omitempty"` // Used when tag is suggested
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TagRelationship represents the relationship between two tags
type TagRelationship struct {
	Tag1ID        int       `json:"tag1_id"`
	Tag2ID        int       `json:"tag2_id"`
	Strength      float64   `json:"strength"`
	CoOccurrences int       `json:"co_occurrences"`
	LastUpdated   time.Time `json:"last_updated"`
}

// TagSuggestion represents a suggested tag with its confidence score
type TagSuggestion struct {
	Tag        *Tag    `json:"tag"`
	Confidence float64 `json:"confidence"`
	Source     string  `json:"source"` // "vector" or "graph"
}

// NewTag creates a new Tag instance
func NewTag(name string, category string, vector []float64) *Tag {
	now := time.Now()
	return &Tag{
		Name:       name,
		Category:   category,
		Vector:     vector,
		UsageCount: 0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// UpdateVector updates the tag's vector representation
func (t *Tag) UpdateVector(vector []float64) {
	t.Vector = vector
	t.UpdatedAt = time.Now()
}

// IncrementUsage increments the tag's usage count
func (t *Tag) IncrementUsage() {
	t.UsageCount++
	t.UpdatedAt = time.Now()
}

// SetConfidence sets the confidence score for a tag suggestion
func (t *Tag) SetConfidence(confidence float64) {
	t.Confidence = confidence
}
