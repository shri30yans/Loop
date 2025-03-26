package utils

import (
	"Loop_backend/internal/models"
	"Loop_backend/platform/database/neo4j/entities"
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

// GetProjectAnalysisPrompt returns the prompt for project analysis
func GetProjectAnalysisPrompt(project *models.Project) string {
	var combinedSections strings.Builder
	for _, section := range project.Sections {
		combinedSections.WriteString(section.Content)
		combinedSections.WriteString("\n")
	}
	text := combinedSections.String()

	return `---Goal---
        Given a text document that is potentially relevant to this activity and a list of entity types, identify all entities of those types from the text and all relationships among the identified entities.

        ---Steps---
        1. Identify all entities. For each identified entity, extract the following information:
        - entity_name: Name of the entity, use same language as input text. If English, capitalized the name.
        - entity_type: One of the following types: [` + DefaultEntityTypes + `]
        - entity_description: Comprehensive description of the entity's attributes and activities
        Format each entity as ("entity"` + TupleDelimiter + `<entity_name>` + TupleDelimiter + `<entity_type>` + TupleDelimiter + `<entity_description>)

        2. From the entities identified in step 1, identify all pairs of (source_entity, target_entity) that are *clearly related* to each other.
        For each pair of related entities, extract the following information:
        - source_entity: name of the source entity, as identified in step 1
        - target_entity: name of the target entity, as identified in step 1
        - relationship_description: explanation as to why you think the source entity and the target entity are related to each other
        - relationship_strength: a numeric score indicating strength of the relationship between the source entity and target entity
        - relationship_keywords: one or more high-level key words that summarize the overarching nature of the relationship
        Format each relationship as ("relationship"` + TupleDelimiter + `<source_entity>` + TupleDelimiter + `<target_entity>` + TupleDelimiter + `<relationship_description>` + TupleDelimiter + `<relationship_keywords>` + TupleDelimiter + `<relationship_strength>)

        3. Identify high-level key words that summarize the main concepts, themes, or topics of the entire text.
        Format the content-level key words as ("content_keywords"` + TupleDelimiter + `<high_level_keywords>)

        4. Return output in English as a single list of all the entities and relationships. Use **` + RecordDelimiter + `** as the list delimiter.

        5. When finished, output ` + CompletionDelimiter + `

        Text:
        ` + text
}
