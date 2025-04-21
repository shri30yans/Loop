package utils

import (
	"Loop_backend/platform/database/neo4j/entities"
	"fmt"
	"strings"
)

// GetCypherTransformPrompt returns a prompt for transforming natural language to Cypher queries
func GetCypherTransformPrompt(query string) string {
	fmt.Println("Input query:", query)
	entityTypes := entities.EntityTypes()

	return fmt.Sprintf(`---Goal---
You are a semantic search engine that converts natural language queries to precise Neo4j Cypher queries.

---Database Schema---
Node Types: %s

Project properties:
- id: Unique project identifier (UUID)
- name/title: Project name (string)
- description: Short project description (string)
- introduction: Detailed project introduction (text)
- status: Project status (draft, published, completed, archived)
- owner_id: UUID of the user who created the project
- created_at: Timestamp when project was created
- updated_at: Timestamp when project was last updated
- tags: Array of tag strings associated with the project

---Relationships---
- RELATED_TO: Generic relationship between entities
- USES: Project uses methodology or technology
- IMPLEMENTS: Project implements a feature
- BELONGS_TO: Project belongs to a category
- HAS_STAKEHOLDER: Project has a stakeholder
- DEVELOPED_BY: Project developed by a person
- DEVELOPMENT: Project has development-related technology

---Case Sensitivity---
Neo4j is case sensitive. Always use these patterns for case-insensitive matching:
- For exact matches: WHERE toLower(node.property) = toLower("value")
- For partial matches: WHERE toLower(node.property) CONTAINS toLower("value")
- Incorrect pattern (don't use): {name: toLower("value")}

---Semantic Understanding---
For ownership/creation queries:
- "Who owns/created X" → Look for (Person)-[:DEVELOPED_BY]-(Project)
- "Find person responsible for X" → Person connected to Project

For similarity queries:
- "Projects similar to X" → Projects sharing tags/categories/technologies
- "Find X related to Y" → Look for nodes with connecting relationships
---Date/Time Operations---
- For date calculations, use duration: date() - duration('P1Y') for 1 year
- For month calculations: date() - duration('P6M') for 6 months
- For day calculations: date() - duration('P30D') for 30 days

---Query Examples---

Example 1 (Project details with relationships):
MATCH (p:Project)
WHERE toLower(p.name) = toLower("ProjectName")
OPTIONAL MATCH (p)-[:USES]->(t:Technology)
OPTIONAL MATCH (p)-[:IMPLEMENTS]->(f:Feature)
OPTIONAL MATCH (p)-[:BELONGS_TO]->(c:Category)
OPTIONAL MATCH (p)-[:HAS_STAKEHOLDER]->(s:Stakeholder)
OPTIONAL MATCH (p)-[:DEVELOPED_BY]->(person:Person)
RETURN p.id, p.name, p.description,
       collect(DISTINCT t.name) AS technologies,
       collect(DISTINCT f.name) AS features,
       collect(DISTINCT c.name) AS categories,
       collect(DISTINCT s.name) AS stakeholders,
       collect(DISTINCT person.name) AS developers

Example 2 (Finding project owner):
MATCH (p:Project)
WHERE toLower(p.name) = toLower("ProjectName")
MATCH (person:Person)-[:DEVELOPED_BY]-(p)
RETURN person.name AS owner, p.name AS project

Example 3 (Similar projects):
MATCH (p1:Project)
WHERE toLower(p1.name) = toLower("ProjectName")
MATCH (p1)-[:BELONGS_TO]->(c:Category)<-[:BELONGS_TO]-(p2:Project)
WHERE p1 <> p2
RETURN p2.name AS similar_project, c.name AS shared_category

---Instructions---
1. Analyze the query for semantic intent (ownership, similarity, features, etc.)
2. Always use toLower() for case-insensitive matching on both sides
3. Use OPTIONAL MATCH for related entities that may not exist
4. Return results with meaningful labels using AS
5. For collection results, use collect(DISTINCT x.property)
6. For time-based queries, use proper duration syntax (e.g., duration('P1Y'))
7. Structure queries to return exactly what was asked for
8. Return ONLY the Cypher query, no explanations

User Query: "%s"
`, strings.Join(entityTypes, ", "), query)
}
