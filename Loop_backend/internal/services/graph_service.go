package services

import (
	"Loop_backend/internal/models"
	"Loop_backend/internal/repositories"
	"Loop_backend/internal/utils"
	"Loop_backend/platform/database/neo4j/entities"
	"fmt"

	"github.com/google/uuid"
)

// GraphService defines operations for managing graph data
type GraphService interface {
	StoreProjectGraph(project *models.Project, graph *models.KnowledgeGraph) error
	GetProjectGraph(projectID uuid.UUID) (*models.KnowledgeGraph, error)
}

// DefaultGraphService implements GraphService
type DefaultGraphService struct {
	graphRepo           repositories.GraphRepository
	relationshipManager *utils.RelationshipManager
}

// NewGraphService creates a new graph service instance
func NewGraphService(graphRepo repositories.GraphRepository) (GraphService, error) {
	relationshipManager, err := utils.GetRelationshipManager()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize relationship manager: %w", err)
	}

	return &DefaultGraphService{
		graphRepo:           graphRepo,
		relationshipManager: relationshipManager,
	}, nil
}


// TODO: rewrite this functiom
// StoreProjectGraph stores a project's knowledge graph
func (s *DefaultGraphService) StoreProjectGraph(project *models.Project, graph *models.KnowledgeGraph) error {
	// Ensure we have a project node
	var projectEntity *models.Entity
	for _, entity := range graph.Entities {
		if entity.Type == "Project" {
			projectEntity = &entity
			break
		}
	}

	// If no project entity exists in the graph, create one from the Project model
	if projectEntity == nil {
		projectEntity = &models.Entity{
			Name:        project.Title,
			Type:        "Project",
			Description: project.Description,
			Properties: map[string]interface{}{
				"project_id": project.ProjectID.String(),
			},
		}
		graph.Entities = append(graph.Entities, *projectEntity)
	}

	// Add owner relationship
	ownerEntity := models.Entity{
		Name:        project.OwnerID.String(),
		Type:        "Person",
		Description: "Project owner",
	}
	graph.Entities = append(graph.Entities, ownerEntity)

	// Add DEVELOPED_BY relationship between project and owner
	ownerRelationship := models.Relationship{
		Source:      projectEntity.Name,
		Target:      ownerEntity.Name,
		Type:        entities.DevelopedBy,
		Description: "Project owner",
		Weight:      10,
		Category:    "ownership",
	}
	graph.Relationships = append(graph.Relationships, ownerRelationship)

	// Add keyword entities and relationships
	for _, keyword := range graph.Keywords {
		keywordEntity := models.Entity{
			Name:        keyword,
			Type:        entities.TypeTag,
			Description: "Project keyword",
			Properties: map[string]interface{}{
				"project_id": project.ProjectID.String(),
			},
		}
		graph.Entities = append(graph.Entities, keywordEntity)

		// Add RELATED_TO relationship between project and keyword
		keywordRelationship := models.Relationship{
			Source:      projectEntity.Name,
			Target:      keyword,
			Type:        entities.RelatedTo,
			Description: "Project keyword",
			Weight:      5,
			Category:    "tags",
		}
		graph.Relationships = append(graph.Relationships, keywordRelationship)
	}

	// Add project_id to relevant entities
	for i := range graph.Entities {
		switch graph.Entities[i].Type {
		case entities.TypeTechnology, entities.TypeFeature, entities.TypeTag:
			if graph.Entities[i].Properties == nil {
				graph.Entities[i].Properties = make(map[string]interface{})
			}
			graph.Entities[i].Properties["project_id"] = project.ProjectID.String()
		}
	}

	// Validate and organize entities by type
	entityMap := make(map[string][]models.Entity)
	for _, entity := range graph.Entities {
		if !s.relationshipManager.ValidateEntityType(entity.Type) {
			return fmt.Errorf("invalid entity type: %s", entity.Type)
		}
		entityMap[entity.Type] = append(entityMap[entity.Type], entity)
	}

	// Process relationships with proper types and additional attributes
	processedRelationships := make([]models.Relationship, 0, len(graph.Relationships))
	for _, rel := range graph.Relationships {
		// Find source and target entity types
		var sourceType, targetType string
		for _, entity := range graph.Entities {
			if entity.Name == rel.Source {
				sourceType = entity.Type
			}
			if entity.Name == rel.Target {
				targetType = entity.Type
			}
			if sourceType != "" && targetType != "" {
				break
			}
		}

		if sourceType == "" || targetType == "" {
			return fmt.Errorf("could not find entity types for relationship: %s -> %s",
				rel.Source, rel.Target)
		}

		// Use provided relationship type or default to RELATED_TO
		relationType := entities.GetRelationshipType(sourceType, targetType)

		// Create processed relationship with additional attributes
		processedRel := models.Relationship{
			Source:      rel.Source,
			Target:      rel.Target,
			Type:        relationType,
			Description: rel.Description,
			Weight:      rel.Weight,
			Category:    rel.Category,
		}
		processedRelationships = append(processedRelationships, processedRel)
	}

	// Update the graph with processed relationships
	graph.Relationships = processedRelationships

	return s.graphRepo.StoreProjectGraph(project.ProjectID, project.OwnerID.String(), graph)
}

// GetProjectGraph retrieves a project's knowledge graph
func (s *DefaultGraphService) GetProjectGraph(projectID uuid.UUID) (*models.KnowledgeGraph, error) {
	return s.graphRepo.GetProjectGraph(projectID)
}
