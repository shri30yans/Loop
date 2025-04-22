	package utils

import (
	"Loop_backend/internal/models"
	"Loop_backend/platform/database/neo4j/entities"
	"fmt"
	"strings"
)

const (
	// Delimiters used in prompts and response parsing
	TupleDelimiter      = ":"       // Separates fields within a record
	RecordDelimiter     = "\n---\n" // Separates individual records in the response
	CompletionDelimiter = "[END]"   // Marks the end of the LLM response
)

var (
	DefaultEntityTypes = strings.Join(entities.EntityTypes(), ", ")
)

// ExtractedEntity represents an entity node extracted from LLM responses
type ExtractedEntity struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Keywords    []string               `json:"keywords,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
}

// ExtractedRelationship represents a relationship between two entities extracted from LLM responses
type ExtractedRelationship struct {
	SourceEntity  string                 `json:"source_entity"`
	TargetEntity  string                 `json:"target_entity"`
	Type          string                 `json:"type"`
	Description   string                 `json:"description"`
	Keywords      []string               `json:"keywords,omitempty"`
	StrengthScore int                    `json:"strength_score"`
	Properties    map[string]interface{} `json:"properties,omitempty"`
}

// GetProjectAnalysisPrompt returns the prompt for project analysis with dynamic relationship typing
func GetProjectAnalysisPrompt(project *models.Project) string {
	var combinedSections strings.Builder
	for _, section := range project.Sections {
		combinedSections.WriteString(section.Content)
		combinedSections.WriteString("\n")
	}
	text := combinedSections.String()
	
	// Generate entity type explanations for the prompt
	entityTypeExplanations := generateEntityTypeExplanations()
	
	return `---Goal---
        Given a text document that is potentially relevant to this activity and a list of entity types, identify all entities of those types from the text and all relationships among the identified entities.
        
        ---Entity Types---
        ` + entityTypeExplanations + `
        
        ---Steps---
        1. Identify all entities. For each identified entity, extract the following information:
        - entity_name: Name of the entity, use same language as input text. If English, capitalize the name.
        - entity_type: One of the following types: [` + DefaultEntityTypes + `]
        - entity_description: Comprehensive description of the entity's attributes and activities
        Format each entity as ("entity"` + TupleDelimiter + `<entity_name>` + TupleDelimiter + `<entity_type>` + TupleDelimiter + `<entity_description>)
        
        2. From the entities identified in step 1, identify all pairs of (source_entity, target_entity) that are *clearly related* to each other.
        For each pair of related entities, extract the following information:
        - source_entity: name of the source entity, as identified in step 1
        - target_entity: name of the target entity, as identified in step 1
        - relationship_description: explanation as to why you think the source entity and the target entity are related to each other
        - relationship_strength: a numeric score indicating strength of the relationship between the source entity and target entity (1-10)
        - relationship_keywords: one or more high-level key words that summarize the overarching nature of the relationship
        Format each relationship as ("relationship"` + TupleDelimiter + `<source_entity>` + TupleDelimiter + `<target_entity>` + TupleDelimiter + `<relationship_description>` + TupleDelimiter + `<relationship_keywords>` + TupleDelimiter + `<relationship_strength>)
        
        3. Identify high-level key words that summarize the main concepts, themes, or topics of the entire text.
        Format the content-level key words as ("content_keywords"` + TupleDelimiter + `<high_level_keywords>)
        
        4. Return output in English as a single list of all the entities and relationships. Use **` + RecordDelimiter + `** as the list delimiter.
        
        5. When finished, output ` + CompletionDelimiter + `
        
        Text:
        ` + text
}

// GetSummaryPrompt generates a prompt for summarizing the project with focus on technologies
func GetSummaryPrompt(project *models.Project) string {
	var combinedSections strings.Builder
	for _, section := range project.Sections {
		combinedSections.WriteString(section.Content)
		combinedSections.WriteString("\n")
	}
	text := combinedSections.String()

	return `---Goal---
		Given a text document about a project, extract and summarize key elements focusing on technologies and other important entities.
		
		---Steps---
		1. Summarize the overall project in 2-3 sentences.
		
		2. Extract all entities by type with the following details for each:
		- Entity Type: One of [` + DefaultEntityTypes + `]
		- Entity Name: The name of the entity
		- Description: Brief description of what it is and its role in the project
		- Relationship Strength: On a scale of 1-10, how important is this entity to the project
		
		Format as tables by entity type, for example:
		
		## Technologies
		| Technology Name | Description | Relationship Strength |
		|----------------|-------------|----------------------|
		| React | Frontend JavaScript library used for UI components | 9 |
		| Node.js | Backend JavaScript runtime environment | 8 |
		
		## Stakeholders
		| Stakeholder Name | Description | Relationship Strength |
		|-----------------|-------------|----------------------|
		| Marketing Team | Responsible for product promotion | 7 |
		
		3. Include tables for all relevant entity types present in the document.
		
		4. End with 3-5 key insights or recommendations based on your analysis.
		
		Text:
		` + text
}

// generateEntityTypeExplanations creates descriptions of entity types for the prompt
func generateEntityTypeExplanations() string {
	entityExplanations := map[string]string{
		entities.TypeProject:     "A project, initiative, or undertaking with defined objectives",
		entities.TypeTechnology:  "Software, tools, programming languages, frameworks or technical components",
		entities.TypeFeature:     "Specific functionality or capability of a system",
		entities.TypeTag:         "Labels or markers used to classify or group items",
		entities.TypeCategory:    "Classification or grouping system",
		entities.TypeStakeholder: "People, roles or groups who have interest or influence in the project",
		entities.TypePerson:      "Specific individual mentioned in the content",
		entities.TypePlatform:    "Underlying systems where software runs or is deployed",
		entities.TypeMethodology: "Approach, framework or set of methods used in development",
	}

	var explanationsBuilder strings.Builder
	for entityType, explanation := range entityExplanations {
		explanationsBuilder.WriteString(fmt.Sprintf("- %s: %s\n", entityType, explanation))
	}
	
	return explanationsBuilder.String()
}

// ProcessEntityRelationships parses the LLM response and generates properly typed relationships
func ProcessEntityRelationships(llmResponse string) ([]ExtractedEntity, []ExtractedRelationship, []string, error) {
	// Remove the completion delimiter if present
	llmResponse = strings.TrimSuffix(llmResponse, CompletionDelimiter)

	// Split response by record delimiter
	records := strings.Split(llmResponse, RecordDelimiter)

	var extractedEntities []ExtractedEntity
	var extractedRelationships []ExtractedRelationship
	var contentKeywords []string

	// Map to store entities by name for quick lookup (case-insensitive)
	entityMap := make(map[string]ExtractedEntity)
	
	// Make sure we have a Project entity
	projectFound := false

	// Process each record
	for _, record := range records {
		record = strings.TrimSpace(record)
		if record == "" {
			continue
		}

		// Entity
		if strings.HasPrefix(record, "(\"entity\"") {
			entityRecord := strings.TrimPrefix(record, "(\"entity\""+TupleDelimiter)
			entityRecord = strings.TrimSuffix(entityRecord, ")")
			entityParts := strings.Split(entityRecord, TupleDelimiter)

			if len(entityParts) >= 3 {
				name := strings.TrimSpace(entityParts[0])
				rawType := strings.TrimSpace(entityParts[1])
				
				// Properly normalize the entity type
				normalizedType := normalizeEntityType(rawType)
				
				// Check if we found a Project entity
				if normalizedType == entities.TypeProject {
					projectFound = true
				}

				entity := ExtractedEntity{
					Name:        name,
					Type:        normalizedType,
					Description: entityParts[2],
				}

				extractedEntities = append(extractedEntities, entity)
				entityMap[strings.ToLower(name)] = entity // Case-insensitive key
			}
		}

		// Relationship
		if strings.HasPrefix(record, "(\"relationship\"") {
			relationshipRecord := strings.TrimPrefix(record, "(\"relationship\""+TupleDelimiter)
			relationshipRecord = strings.TrimSuffix(relationshipRecord, ")")
			relParts := strings.Split(relationshipRecord, TupleDelimiter)

			if len(relParts) >= 5 {
				sourceName := strings.TrimSpace(relParts[0])
				targetName := strings.TrimSpace(relParts[1])
				lookupSource := strings.ToLower(sourceName)
				lookupTarget := strings.ToLower(targetName)

				sourceEntity, sourceExists := entityMap[lookupSource]
				targetEntity, targetExists := entityMap[lookupTarget]

				if sourceExists && targetExists {
					var keywords []string
					if len(relParts) > 3 && relParts[3] != "" {
						for _, keyword := range strings.Split(relParts[3], ",") {
							keywords = append(keywords, strings.TrimSpace(keyword))
						}
					}

					// Correctly determine relationship type using the entity types
					relType := entities.GetRelationshipType(sourceEntity.Type, targetEntity.Type)

					relationship := ExtractedRelationship{
						SourceEntity:  sourceEntity.Name,
						TargetEntity:  targetEntity.Name,
						Type:          relType,
						Description:   relParts[2],
						Keywords:      keywords,
						StrengthScore: parseStrengthScore(relParts[4]),
					}

					extractedRelationships = append(extractedRelationships, relationship)
				}
			}
		}

		// Content-level keywords
		if strings.HasPrefix(record, "(\"content_keywords\"") {
			keywordsRecord := strings.TrimPrefix(record, "(\"content_keywords\""+TupleDelimiter)
			keywordsRecord = strings.TrimSuffix(keywordsRecord, ")")

			for _, keyword := range strings.Split(keywordsRecord, ",") {
				contentKeywords = append(contentKeywords, strings.TrimSpace(keyword))
			}
		}
	}
	
	// If we have project data but no explicit project entity, add one
	if !projectFound && len(extractedEntities) > 0 {
		// Create a default project entity based on the first section title or "Project"
		projectEntity := ExtractedEntity{
			Name:        "AgroLink", // Default project name - you could make this dynamic
			Type:        entities.TypeProject,
			Description: "Main project from analysis",
		}
		extractedEntities = append(extractedEntities, projectEntity)
		entityMap[strings.ToLower(projectEntity.Name)] = projectEntity
		
		// Link other entities to the project with appropriate relationships
		for _, entity := range extractedEntities {
			if entity.Name != projectEntity.Name {
				// Determine relationship type based on entity type
				relType := entities.GetRelationshipType(projectEntity.Type, entity.Type)
				
				// Create relationship
				relationship := ExtractedRelationship{
					SourceEntity:  projectEntity.Name,
					TargetEntity:  entity.Name,
					Type:          relType,
					Description:   fmt.Sprintf("Auto-generated relationship between project and %s", entity.Type),
					Keywords:      []string{"auto-generated"},
					StrengthScore: 7, // Default medium-high strength
				}
				extractedRelationships = append(extractedRelationships, relationship)
			}
		}
	}

	return extractedEntities, extractedRelationships, contentKeywords, nil
}


// determineRelationshipType uses the relationship patterns from entities package
func determineRelationshipType(sourceType, targetType string) string {
    // Use the GetRelationshipType function from the entities package
    return entities.GetRelationshipType(sourceType, targetType)
}

// normalizeEntityType ensures we're using the proper case for entity types
func normalizeEntityType(entityType string) string {
	// First trim any whitespace
	entityType = strings.TrimSpace(entityType)
	
	// Check for direct match with existing types (case-insensitive)
	for _, validType := range entities.EntityTypes() {
		if strings.EqualFold(entityType, validType) {
			return validType // Return the properly cased version
		}
	}
	
	// Try to normalize the type
	// Some special case normalizations
	switch strings.ToLower(entityType) {
	case "tech", "technology", "technologies":
		return entities.TypeTechnology
	case "person", "individual", "people", "user":
		return entities.TypePerson
	case "stakeholder", "stakeholders", "client", "customer":
		return entities.TypeStakeholder
	case "project", "application", "app", "system", "product":
		return entities.TypeProject
	case "feature", "functionality", "function", "capability":
		return entities.TypeFeature
	case "tag", "label", "keyword":
		return entities.TypeTag
	case "category", "group", "classification":
		return entities.TypeCategory
	case "methodology", "method", "framework", "approach":
		return entities.TypeMethodology
	case "platform", "infrastructure", "server", "service":
		return entities.TypePlatform
	}
	
	// If we still can't determine the type, default to Tag
	return entities.TypeTag
}

// parseStrengthScore converts string strength score to int with error handling
func parseStrengthScore(scoreStr string) int {
	scoreStr = strings.TrimSpace(scoreStr)
	score := 0
	_, err := fmt.Sscanf(scoreStr, "%d", &score)
	if err != nil {
		// Default to 5 (medium strength) if parsing fails
		return 5
	}
	// Ensure score is within bounds
	if score < 1 {
		return 1
	}
	if score > 10 {
		return 10
	}
	return score
}

// ConvertToEntitiesFormat converts ExtractedEntity to entities.Entity
func ConvertToEntitiesFormat(extracted []ExtractedEntity) []models.Entity {
	result := make([]models.Entity, 0, len(extracted))
	for _, e := range extracted {
		result = append(result, models.Entity{
			Name:        e.Name,
			Type:        e.Type,
			Description: e.Description,
		})
	}
	return result
}

// ConvertToRelationshipsFormat converts ExtractedRelationship to models.Relationship
// ConvertToRelationshipsFormat converts ExtractedRelationship to models.Relationship
func ConvertToRelationshipsFormat(extracted []ExtractedRelationship) []models.Relationship {
    result := make([]models.Relationship, 0, len(extracted))
    for _, r := range extracted {
        result = append(result, models.Relationship{
            Source:      r.SourceEntity,      // Changed from SourceEntity
            Target:      r.TargetEntity,      // Changed from TargetEntity
            Type:        r.Type,              // This field name seems correct
            Description: r.Description,       // This field name seems correct

        })
    }
    return result
}
