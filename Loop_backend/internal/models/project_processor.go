package models

import (
	"github.com/google/uuid"
)

// ProjectProcessor defines the interface for project analysis
type ProjectProcessor interface {
	AnalyzeNewProject(project *Project) error
}

// KnowledgeGraphService defines operations for knowledge graphs
type KnowledgeGraphService interface {
	StoreProjectGraph(projectID uuid.UUID, graph *KnowledgeGraph) error
	GetProjectGraph(projectID uuid.UUID) (*KnowledgeGraph, error)
}

// EntityProcessor defines the interface for entity processing
type EntityProcessor interface {
	ExtractEntities(text string) ([]Entity, error)
	ExtractRelationships(text string) ([]Relationship, error)
	GenerateKnowledgeGraph(text string) (*KnowledgeGraph, error)
	ConvertToTags(entities []Entity, relationships []Relationship) []Tag
}
