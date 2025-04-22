package impl

import (
    "Loop_backend/internal/models"
    "Loop_backend/internal/repositories/graph/queries"

    "github.com/google/uuid"
    "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// GraphRepository defines all graph database operations
type GraphRepository interface {
    EntityRepository
    ProjectRepository
    RelationshipRepository
    PersonRepository
}

// EntityRepository handles entity-specific operations
type EntityRepository interface {
    GetProjectEntities(projectID string, entityType string) ([]models.Entity, error)
    GetEntitiesByType(entityType string) ([]models.Entity, error)
    CreateEntity(entity *models.Entity, projectID string) error
    UpdateEntity(entity *models.Entity) error
    DeleteEntity(entityType, entityID string) error
}

// ProjectRepository handles project-specific operations
type ProjectRepository interface {
    StoreProjectGraph(projectID uuid.UUID, ownerID string, graph *models.KnowledgeGraph) error
    GetProjectGraph(projectID uuid.UUID) (*models.KnowledgeGraph, error)
    DeleteProjectNode(projectID string) error
}

// RelationshipRepository handles relationship operations
type RelationshipRepository interface {
    CreateRelationship(source, target string, relType string, props map[string]interface{}) error
    GetProjectRelationships(projectID string) ([]models.Relationship, error)
    GetEntityRelationships(entityType, entityID string) ([]models.Relationship, error)
    UpdateRelationship(relationshipID string, props map[string]interface{}) error
    DeleteRelationship(relationshipID string) error
}

// PersonRepository handles person/user specific operations
type PersonRepository interface {
    GetPersonProjects(personID string) ([]models.Project, error)
    GetPersonContributions(personID string, limit int) ([]models.Contribution, error)
    LinkPersonToEntity(personID, entityID string, entityType string, relationType string, props map[string]interface{}) error
}

// baseRepo provides common functionality for all repository implementations
type baseRepo struct {
    driver       neo4j.Driver
    queryManager *queries.QueryManager
}

// NewBaseRepo creates a new base repository with the given driver
func NewBaseRepo(driver neo4j.Driver) baseRepo {
    return baseRepo{
        driver:       driver,
        queryManager: queries.GetQueryManager(),
    }
}

// NewEntityRepo creates a new entity repository
func NewEntityRepo(base baseRepo) EntityRepository {
    return &entityRepo{baseRepo: base}
}

// NewProjectRepo creates a new project repository
func NewProjectRepo(base baseRepo) ProjectRepository {
    return &projectRepo{baseRepo: base}
}

// NewRelationshipRepo creates a new relationship repository
func NewRelationshipRepo(base baseRepo) RelationshipRepository {
    return &relationshipRepo{baseRepo: base}
}

// NewPersonRepo creates a new person repository
func NewPersonRepo(base baseRepo) PersonRepository {
    return &personRepo{baseRepo: base}
}
