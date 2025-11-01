package graph

import (
	"Loop_backend/internal/dto"
	"context"
)

// GraphService defines the interface for graph operations
type GraphService interface {
	// Graph operations
	GetProjectsMetadata(ctx context.Context, projectIDs []string) (map[string]ProjectMetadata, error)
	StoreProjectGraph(ctx context.Context, projectID string, graph *KnowledgeGraph) error
	GetProjectGraph(ctx context.Context, projectID string) (*KnowledgeGraph, error)

	// Entity operations
	CreateEntity(ctx context.Context, entity *Entity) error
	GetEntityByName(ctx context.Context, name string) (*Entity, error)
	GetEntitiesByType(ctx context.Context, entityType string) ([]Entity, error)

	// Relationship operations
	CreateRelationship(ctx context.Context, source, target, relType string, weight int) error
	GetRelationships(ctx context.Context, entityName string) ([]Relationship, error)
	GetRelatedProjects(ctx context.Context, entityName string) ([]string, error)

	// Graph processing
	ProcessGraph(ctx context.Context, graph *KnowledgeGraph) error

	// Search operations
	SearchProjects(ctx context.Context, query string) ([]dto.SearchToolsResult, error)
	SearchByTechnology(ctx context.Context, tech string) ([]dto.SearchToolsResult, error)
	SearchByDomain(ctx context.Context, domain string) ([]dto.SearchToolsResult, error)
	TransformQueryToCypher(query string) (string, error)
}

// ProjectMetadata represents metadata about a project
type ProjectMetadata struct {
	ID           string
	ProjectID    string
	Name         string
	Description  string
	Technologies []string
	Domains      []string
	Features     []string
	Metadata     map[string]interface{}
}

// KnowledgeGraph represents a graph of project knowledge
type KnowledgeGraph struct {
	Name          string
	Description   string
	Entities      []Entity
	Relationships []Relationship
	Metadata      map[string]interface{}
}

// Entity represents a node in the knowledge graph
type Entity struct {
	Name     string
	Type     string
	Level    string
	Metadata map[string]interface{}
}

// Relationship represents an edge in the knowledge graph
type Relationship struct {
	Source   string
	Target   string
	Type     string
	Weight   int
	Metadata map[string]interface{}
}
