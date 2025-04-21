package utils

import (
	"Loop_backend/internal/repositories"
	"Loop_backend/platform/database/neo4j/entities"
	"fmt"
	"strings"
)

// GetTopicSearchPrompt returns a prompt specifically for topic-based search
func GetTopicSearchPrompt(query string, graphRepo repositories.GraphRepository) string {
	fmt.Println("Topic search query:", query)
	entityTypes := entities.EntityTypes()

	// Fetch all relationship types from the database
	relationships, err := graphRepo.GetAllRelationshipTypes()
	if err != nil {
		fmt.Printf("Error fetching relationship types: %v\n", err)
		// Fall back to hardcoded relationships if query fails
		relationships = []string{"HAS_TAG", "BELONGS_TO", "USES", "RELATED_TO", "IMPLEMENTS", "DEVELOPED_BY", "DEVELOPMENT"}
	}

	// Debug print to verify relationships
	fmt.Printf("Found %d relationships: %v\n", len(relationships), relationships)

	// Create relationship strings
	relationshipStrings := make([]string, len(relationships))
	for i, rel := range relationships {
		relationshipStrings[i] = fmt.Sprintf("- %s", rel)
	}

	// Make sure the relationships section appears in the prompt
	relationshipsSection := strings.Join(relationshipStrings, "\n")
	if relationshipsSection == "" {
		relationshipsSection = "- HAS_TAG (connects Projects to Tags)\n- USES (connects Projects to Technologies)"
	}

	return fmt.Sprintf(`---Goal---
You are a semantic search engine that converts natural language topic queries to Neo4j Cypher queries.
The goal is to find projects semantically related to the specified topic, category, or concept.

---Database Schema---
Node Types: %s

---Relationships---
%s

---IMPORTANT: Search Priority---
Projects are frequently categorized by tags, which is the MOST RELIABLE way to find them by topic.
ALWAYS include tag-based searches as your primary strategy, then fallback to other approaches.

---Semantic Search Approach---
For effective semantic search:
1. FIRST check tags (this is highest priority)
2. Check related entities (categories, technologies)
3. Check project properties (name, description)
4. Use partial matching with CONTAINS
5. Consider related concepts and synonyms (agriculture → farming, crops)
6. For multi-word topics, also search for individual words

---Instructions---
1. Extract key concepts from the query
2. Create a search query with TAG MATCHING as the FIRST approach
3. Use case-insensitive matching with toLower() on both sides
4. Return project ID as "projectId" and name as "projectName"
5. Use a UNION approach to combine different search strategies
6. Return ONLY the Cypher query - no explanations

---Examples---

Example 1: Finding projects tagged with related concepts
MATCH (p:Project)-[:HAS_TAG]->(t:Tag)
WHERE toLower(t.name) CONTAINS toLower("agriculture") 
   OR toLower(t.name) CONTAINS toLower("farming")
   OR toLower(t.name) CONTAINS toLower("crops")
RETURN DISTINCT p.id as projectId, p.name as projectName
UNION
MATCH (p:Project)
WHERE toLower(p.name) CONTAINS toLower("agriculture") 
   OR toLower(p.description) CONTAINS toLower("agriculture")
RETURN DISTINCT p.id as projectId, p.name as projectName

Example 2: For multi-word topics like "machine learning"
MATCH (p:Project)-[:HAS_TAG]->(t:Tag)
WHERE toLower(t.name) CONTAINS toLower("machine learning")
   OR (toLower(t.name) CONTAINS toLower("machine") AND toLower(t.name) CONTAINS toLower("learning"))
   OR toLower(t.name) CONTAINS toLower("ml")
   OR toLower(t.name) CONTAINS toLower("artificial intelligence")
RETURN DISTINCT p.id as projectId, p.name as projectName
UNION
MATCH (p:Project)
WHERE toLower(p.name) CONTAINS toLower("machine learning")
   OR toLower(p.description) CONTAINS toLower("machine learning")
RETURN DISTINCT p.id as projectId, p.name as projectName

Example 3: For technical topics like "predictive analysis"
MATCH (p:Project)-[:HAS_TAG]->(t:Tag)
WHERE toLower(t.name) CONTAINS toLower("predictive") 
   OR toLower(t.name) CONTAINS toLower("analysis")
   OR toLower(t.name) CONTAINS toLower("predictive analysis")
   OR toLower(t.name) CONTAINS toLower("data science")
RETURN DISTINCT p.id as projectId, p.name as projectName
UNION
MATCH (p:Project)
WHERE toLower(p.description) CONTAINS toLower("predictive") 
   OR toLower(p.description) CONTAINS toLower("analysis")
RETURN DISTINCT p.id as projectId, p.name as projectName

User Query: "%s"
`, strings.Join(entityTypes, ", "), relationshipsSection, query)
}
