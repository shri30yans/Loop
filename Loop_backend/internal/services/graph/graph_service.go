package graph

import (
    "Loop_backend/internal/ai/interfaces"
    "Loop_backend/internal/models"
    graphRepo "Loop_backend/internal/repositories/graph"
)

// Service provides graph-related operations
type Service interface {
    StoreProjectGraph(projectID string, graph *models.KnowledgeGraph) error
}

// DefaultService implements Service
type DefaultService struct {
    graphRepo graphRepo.Repository
    provider  interfaces.Provider
}

// NewService creates a new graph service
func NewService(repo graphRepo.Repository, provider interfaces.Provider) Service {
    return &DefaultService{
        graphRepo: repo,
        provider:  provider,
    }
}

// StoreProjectGraph stores a project's knowledge graph
func (s *DefaultService) StoreProjectGraph(projectID string, graph *models.KnowledgeGraph) error {
    return s.graphRepo.StoreProjectGraph( projectID, graph)
}
