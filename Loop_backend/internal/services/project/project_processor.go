package project

import (
	"Loop_backend/internal/ai/providers"
	"Loop_backend/internal/models"
	"Loop_backend/internal/services/graph"
	entitiesPackage "Loop_backend/platform/database/neo4j/entities"
	"fmt"
	"regexp"
	"strings"
)

// DefaultProjectProcessor implements models.ProjectProcessor
type DefaultProjectProcessor struct {
	provider     providers.Provider
	graphService graph.GraphService
}

// NewProjectProcessor creates a new project analyzer instance
func NewProjectProcessor(
	provider providers.Provider,
	graphService graph.GraphService,
) models.ProjectProcessor {
	return &DefaultProjectProcessor{
		provider:     provider,
		graphService: graphService,
	}
}

// extractEntities handles the AI-based entity extraction
func (pa *DefaultProjectProcessor) extractEntities(project *models.Project) (*providers.ChatResponse, error) {
	prompt := GetProjectAnalysisPrompt(project)

	response, err := pa.provider.Chat([]providers.Message{
		{Role: "user", Content: prompt},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}

	return response, nil
}

// validateEntities ensures all entities have valid types
func (pa *DefaultProjectProcessor) validateEntities(entities []models.Entity) []models.Entity {
	validEntities := make([]models.Entity, 0, len(entities))

	for _, entity := range entities {
		if entitiesPackage.IsValidEntityType(entity.Type) {
			validEntities = append(validEntities, entity)
		}
	}

	return validEntities
}

// ensureProjectEntity adds project entity if not present
func (pa *DefaultProjectProcessor) ensureProjectEntity(entities []models.Entity, project *models.Project) []models.Entity {
	hasProject := false
	for _, entity := range entities {
		if entity.Type == entitiesPackage.TypeProject {
			hasProject = true
			break
		}
	}

	if !hasProject {
		projectEntity := models.Entity{
			Name:        project.Title,
			Type:        entitiesPackage.TypeProject,
			Description: project.Description,
		}
		entities = append(entities, projectEntity)
	}

	return entities
}

// ensureOwnerEntity adds owner entity and relationship
func (pa *DefaultProjectProcessor) ensureOwnerEntity(entities []models.Entity, project *models.Project) []models.Entity {
	ownerEntity := models.Entity{
		Name:        project.OwnerID.String(),
		Type:        entitiesPackage.TypePerson,
		Description: "Project owner",
	}
	return append(entities, ownerEntity)
}

// processResponse parses AI response into knowledge graph
func (pa *DefaultProjectProcessor) processResponse(response *providers.ChatResponse) (*models.KnowledgeGraph, error) {
	fmt.Printf("\n=== Processing Raw Response ===\n")
	// Remove completion delimiter if present
	content := strings.Split(response.Content, "[END]")[0]

	// Split into entity records
	records := strings.Split(content, "\n")
	entities := make([]models.Entity, 0)

	fmt.Printf("Found %d raw records\n", len(records))

	// Parse each record into an entity
	for i, record := range records {
		record = strings.TrimSpace(record)
		if record == "" {
			continue
		}

		fmt.Printf("\nParsing record %d:\n%s\n", i+1, record)
		if entity, err := parseEntity(record); err == nil {
			entities = append(entities, entity)
			fmt.Printf("Successfully parsed entity: %+v\n", entity)
		} else {
			fmt.Printf("Failed to parse entity: %v\n", err)
		}
	}

	fmt.Printf("\nExtracted %d entities\n", len(entities))
	return &models.KnowledgeGraph{
		Entities: entities,
	}, nil
}

// Parse entity from formatted line with tuple delimiters
func parseEntity(line string) (models.Entity, error) {
	pattern := fmt.Sprintf(`\(([^%s]+)%s([^%s]+)%s([^%s]+)\)`, TupleDelimiter, TupleDelimiter, TupleDelimiter, TupleDelimiter, TupleDelimiter)
	entityPattern := regexp.MustCompile(pattern)

	matches := entityPattern.FindStringSubmatch(line)
	if len(matches) != 4 {
		return models.Entity{}, fmt.Errorf("invalid entity format")
	}

	name := strings.TrimSpace(matches[1])
	entityType := strings.TrimSpace(matches[2])
	description := strings.TrimSpace(matches[3])

	if name == "" || entityType == "" {
		return models.Entity{}, fmt.Errorf("empty name or type not allowed")
	}

	if !entitiesPackage.IsValidEntityType(entityType) {
		return models.Entity{}, fmt.Errorf("invalid entity type: %s", entityType)
	}

	return models.Entity{
		Name:        name,
		Type:        entityType,
		Description: description,
	}, nil
}

// np print entities nicely
func debugPrintEntities(entities []models.Entity) {
	fmt.Println("\n=== Entities ===")
	for _, e := range entities {
		fmt.Printf("Type: %s, Name: %s, Description: %s\n", e.Type, e.Name, e.Description)
	}
	fmt.Println("==============")
}

// Debug function to print AI response
func debugPrintResponse(response *providers.ChatResponse) {
	fmt.Println("\n=== AI Response ===")
	fmt.Printf("Content:\n%s\n", response.Content)
	fmt.Println("=================")
}

// Debug function to print knowledge graph
func debugPrintGraph(graph *models.KnowledgeGraph) {
	fmt.Println("\n=== Knowledge Graph ===")
	fmt.Printf("Number of entities: %d\n", len(graph.Entities))
	debugPrintEntities(graph.Entities)
	fmt.Println("====================")
}

// AnalyzeNewProject performs AI analysis on a new project
func (pa *DefaultProjectProcessor) AnalyzeNewProject(project *models.Project) error {
	fmt.Printf("\n\n=== Processing Project: %s ===\n", project.ProjectID)

	// Extract entities using AI
	response, err := pa.extractEntities(project)
	if err != nil {
		return err
	}
	debugPrintResponse(response)

	// Parse response into knowledge graph
	graph, err := pa.processResponse(response)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	debugPrintGraph(graph)

	// Validate and ensure required entities
	graph.Entities = pa.validateEntities(graph.Entities)
	fmt.Println("\n=== After Validation ===")
	debugPrintEntities(graph.Entities)

	graph.Entities = pa.ensureProjectEntity(graph.Entities, project)
	graph.Entities = pa.ensureOwnerEntity(graph.Entities, project)
	fmt.Println("\n=== Final Graph Before Storage ===")
	debugPrintGraph(graph)

	// Store graph in database
	err = pa.graphService.StoreProjectGraph(project.ProjectID, project.OwnerID.String(), graph)
	if err != nil {
		return fmt.Errorf("failed to store graph: %w", err)
	}

	fmt.Printf("\n=== Completed Processing Project: %s ===\n\n", project.ProjectID)
	return nil
}
