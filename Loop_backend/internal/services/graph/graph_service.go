package graph

import (
	"Loop_backend/internal/models"
	graphRepo "Loop_backend/internal/repositories/graph/impl"
	parser "Loop_backend/internal/services/graph/response_parser"

	"github.com/google/uuid"
)

// GraphService handles operations on the graph database
type GraphService interface {
	// Project operations
	StoreProjectGraph(projectID uuid.UUID, ownerID string, graph *models.KnowledgeGraph) error
	GetProjectGraph(projectID uuid.UUID) (*models.KnowledgeGraph, error)
	DeleteProject(projectID string) error

	// Entity operations
	GetProjectEntities(projectID string, entityType string) ([]models.Entity, error)
	GetEntitiesByType(entityType string) ([]models.Entity, error)

	// Relationship operations
	CreateRelationship(source, target string, relType string, props map[string]interface{}) error
	GetProjectRelationships(projectID string) ([]models.Relationship, error)
	UpdateRelationship(relationshipID string, props map[string]interface{}) error
}

type graphService struct {
	graphRepo graphRepo.GraphRepository
	parser    parser.ResponseParser
}

// NewGraphService creates a new graph service
func NewGraphService(repo graphRepo.GraphRepository, responseParser parser.ResponseParser) GraphService {
	return &graphService{
		graphRepo: repo,
		parser:    responseParser,
	}
}

// StoreProjectGraph stores the project graph in the database
func (s *graphService) StoreProjectGraph(projectID uuid.UUID, ownerID string, graph *models.KnowledgeGraph) error {
	return s.graphRepo.StoreProjectGraph(projectID, ownerID, graph)
}

// GetProjectGraph retrieves a project's knowledge graph
func (s *graphService) GetProjectGraph(projectID uuid.UUID) (*models.KnowledgeGraph, error) {
	return s.graphRepo.GetProjectGraph(projectID)
}

// DeleteProject deletes a project and all its relationships
func (s *graphService) DeleteProject(projectID string) error {
	return s.graphRepo.DeleteProjectNode(projectID)
}

// GetProjectEntities gets all entities of a specific type for a project
func (s *graphService) GetProjectEntities(projectID string, entityType string) ([]models.Entity, error) {
	return s.graphRepo.GetProjectEntities(projectID, entityType)
}

// GetEntitiesByType gets all entities of a specific type
func (s *graphService) GetEntitiesByType(entityType string) ([]models.Entity, error) {
	return s.graphRepo.GetEntitiesByType(entityType)
}

// CreateRelationship creates a relationship between two entities
func (s *graphService) CreateRelationship(source, target string, relType string, props map[string]interface{}) error {
	return s.graphRepo.CreateRelationship(source, target, relType, props)
}

// GetProjectRelationships gets all relationships for a project
func (s *graphService) GetProjectRelationships(projectID string) ([]models.Relationship, error) {
	return s.graphRepo.GetProjectRelationships(projectID)
}

// UpdateRelationship updates the properties of a relationship
func (s *graphService) UpdateRelationship(relationshipID string, props map[string]interface{}) error {
	return s.graphRepo.UpdateRelationship(relationshipID, props)
}
