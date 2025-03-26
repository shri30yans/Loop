package services

import (
    "Loop_backend/internal/models"
    "Loop_backend/internal/repositories"
    "github.com/google/uuid"
)

// TagService defines the interface for tag operations
type TagService interface {
    CreateTagWithEmbedding(tag *models.Tag) error
    GetProjectTags(projectID uuid.UUID) ([]models.Tag, error)
}

// DefaultTagService implements TagService
type DefaultTagService struct {
    tagRepo     TagRepository
    graphRepo   repositories.GraphRepository
}

// NewTagService creates a new tag service instance
func NewTagService(
    tagRepo TagRepository,
    graphRepo repositories.GraphRepository,
) TagService {
    return &DefaultTagService{
        tagRepo:   tagRepo,
        graphRepo: graphRepo,
    }
}

// CreateTagWithEmbedding creates a new tag with its embedding
func (s *DefaultTagService) CreateTagWithEmbedding(tag *models.Tag) error {
    // First create the tag in PostgreSQL
    err := s.tagRepo.CreateTag(tag)
    if err != nil {
        return err
    }

    // Then create the tag node in Neo4j
    err = s.graphRepo.CreateTagNode(tag)
    if err != nil {
        return err
    }

    return nil
}

// GetProjectTags retrieves all tags for a project
func (s *DefaultTagService) GetProjectTags(projectID uuid.UUID) ([]models.Tag, error) {
    return s.tagRepo.GetTagsByProjectID(projectID)
}

// TagRepository defines the interface for tag data operations
type TagRepository interface {
    CreateTag(tag *models.Tag) error
    GetTagsByProjectID(projectID uuid.UUID) ([]models.Tag, error)
}
