package project

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

// GetProjectAnalysisPrompt returns the prompt for project analysis
func GetProjectAnalysisPrompt(project *models.Project) string {
	// Build entity type instructions from metadata
	var typeInstructions strings.Builder
	typeInstructions.WriteString("\nEntity Types to Extract:\n")
	for _, entityType := range entities.EntityTypes() {
		info, _ := entities.GetEntityTypeInfo(entityType)
		typeInstructions.WriteString(fmt.Sprintf("- %s: %s\n  Guidelines: %s\n",
			info.Name,
			info.Description,
			info.ExtractionGuidelines,
		))
	}

	// Combine project sections
	var combinedSections strings.Builder
	for _, section := range project.Sections {
		combinedSections.WriteString(section.Content)
		combinedSections.WriteString("\n")
	}
	text := combinedSections.String()

	return fmt.Sprintf(`---Goal---
Extract and classify entities from the given text according to predefined entity types.
Your task is to identify entities that fit into the specified types, providing clear descriptions
for each identified entity.
 Only send the recognised entities in the format (<entity_name>%s<entity_type>%s<entity_description>)
 Do not give any explaination. Only the entities.

---Entity Type Reference---%s

---Instructions---
1. For each identified entity, extract:
   - entity_name: Name of the entity (use same language as input text, capitalize if English)
   - entity_type: Must be one of the listed types above
   - entity_description: Clear and comprehensive description of the entity

2. Format each entity as:
   (<entity_name>%s<entity_type>%s<entity_description>)

   Example:
   (Node.js%sTechnology%sJavaScript Library)
   (GoLang%sTechnology%sProgramming language)
   (Farmers%sStakeholder%sFarming)

   Strictly follow the above format. 

---Input Text---
%s`,
		typeInstructions.String(),
		TupleDelimiter,
		TupleDelimiter,
		TupleDelimiter,
		TupleDelimiter,
		TupleDelimiter,
		TupleDelimiter,
		TupleDelimiter,
		TupleDelimiter,
		RecordDelimiter,
		CompletionDelimiter,
		text,
	)
}
