package graph

import (
    "Loop_backend/internal/models"
    "Loop_backend/internal/repositories/graph/impl"

    "github.com/google/uuid"
    "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// compositeGraphRepository combines all specialized repositories
type compositeGraphRepository struct {
    projectRepo      impl.ProjectRepository
    entityRepo      impl.EntityRepository
    relationshipRepo impl.RelationshipRepository
    personRepo      impl.PersonRepository
}

// NewGraphRepository creates a new instance of the graph repository
func NewGraphRepository(driver neo4j.Driver) impl.GraphRepository {
    if driver == nil {
        return nil
    }

    base := impl.NewBaseRepo(driver)
    return &compositeGraphRepository{
        projectRepo:      impl.NewProjectRepo(base),
        entityRepo:      impl.NewEntityRepo(base),
        relationshipRepo: impl.NewRelationshipRepo(base),
        personRepo:      impl.NewPersonRepo(base),
    }
}

// StoreProjectGraph stores the entire knowledge graph for a project
func (r *compositeGraphRepository) StoreProjectGraph(projectID uuid.UUID, ownerID string, graph *models.KnowledgeGraph) error {
    return r.projectRepo.StoreProjectGraph(projectID, ownerID, graph)
}

// GetProjectGraph retrieves a project's knowledge graph
func (r *compositeGraphRepository) GetProjectGraph(projectID uuid.UUID) (*models.KnowledgeGraph, error) {
    return r.projectRepo.GetProjectGraph(projectID)
}

// DeleteProjectNode deletes a project and all its relationships
func (r *compositeGraphRepository) DeleteProjectNode(projectID string) error {
    return r.projectRepo.DeleteProjectNode(projectID)
}

// GetProjectEntities retrieves all entities of a specific type for a project
func (r *compositeGraphRepository) GetProjectEntities(projectID string, entityType string) ([]models.Entity, error) {
    return r.entityRepo.GetProjectEntities(projectID, entityType)
}

// GetEntitiesByType retrieves all entities of a specific type
func (r *compositeGraphRepository) GetEntitiesByType(entityType string) ([]models.Entity, error) {
    return r.entityRepo.GetEntitiesByType(entityType)
}

// CreateEntity creates a new entity in the graph
func (r *compositeGraphRepository) CreateEntity(entity *models.Entity, projectID string) error {
    return r.entityRepo.CreateEntity(entity, projectID)
}

// UpdateEntity updates an existing entity in the graph
func (r *compositeGraphRepository) UpdateEntity(entity *models.Entity) error {
    return r.entityRepo.UpdateEntity(entity)
}

// DeleteEntity deletes an entity from the graph
func (r *compositeGraphRepository) DeleteEntity(entityType, entityID string) error {
    return r.entityRepo.DeleteEntity(entityType, entityID)
}

// CreateRelationship creates a relationship between two entities
func (r *compositeGraphRepository) CreateRelationship(source, target string, relType string, props map[string]interface{}) error {
    return r.relationshipRepo.CreateRelationship(source, target, relType, props)
}

// GetProjectRelationships retrieves all relationships for a project
func (r *compositeGraphRepository) GetProjectRelationships(projectID string) ([]models.Relationship, error) {
    return r.relationshipRepo.GetProjectRelationships(projectID)
}

// GetEntityRelationships retrieves all relationships for an entity
func (r *compositeGraphRepository) GetEntityRelationships(entityType, entityID string) ([]models.Relationship, error) {
    return r.relationshipRepo.GetEntityRelationships(entityType, entityID)
}

// UpdateRelationship updates the properties of a relationship
func (r *compositeGraphRepository) UpdateRelationship(relationshipID string, props map[string]interface{}) error {
    return r.relationshipRepo.UpdateRelationship(relationshipID, props)
}

// DeleteRelationship deletes a relationship by ID
func (r *compositeGraphRepository) DeleteRelationship(relationshipID string) error {
    return r.relationshipRepo.DeleteRelationship(relationshipID)
}

// GetPersonProjects retrieves all projects associated with a person
func (r *compositeGraphRepository) GetPersonProjects(personID string) ([]models.Project, error) {
    return r.personRepo.GetPersonProjects(personID)
}

// GetPersonContributions retrieves all contributions made by a person
func (r *compositeGraphRepository) GetPersonContributions(personID string, limit int) ([]models.Contribution, error) {
    return r.personRepo.GetPersonContributions(personID, limit)
}

// LinkPersonToEntity creates a relationship between a person and an entity
func (r *compositeGraphRepository) LinkPersonToEntity(personID, entityID string, entityType string, relationType string, props map[string]interface{}) error {
    return r.personRepo.LinkPersonToEntity(personID, entityID, entityType, relationType, props)
}
