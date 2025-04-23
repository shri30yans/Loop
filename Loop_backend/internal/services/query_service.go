package services

import (
	"Loop_backend/internal/ai/interfaces"
	"Loop_backend/internal/repositories"
	"Loop_backend/internal/utils"
	"fmt"
	"strings"
)

// QueryService defines operations for transforming and executing queries
type QueryService interface {
	TransformQueryToCypher(query string) (string, error)
	TransformQueryToTopicCypher(query string) (string, error)
	ExecuteSearchQuery(query string) ([]map[string]interface{}, error)
	ExecuteTopicSearchQuery(query string) ([]map[string]interface{}, error)
	GetAllProjects() ([]map[string]interface{}, error) // Add this
}

// DefaultQueryService implements QueryService
type DefaultQueryService struct {
	provider  interfaces.Provider
	graphRepo repositories.GraphRepository
}

// NewQueryService creates a new query service instance
func NewQueryService(provider interfaces.Provider, graphRepo repositories.GraphRepository) QueryService {
	return &DefaultQueryService{
		provider:  provider,
		graphRepo: graphRepo,
	}
}

// TransformQueryToCypher transforms a natural language query to a Cypher query
func (s *DefaultQueryService) TransformQueryToCypher(query string) (string, error) {
	prompt := utils.GetCypherTransformPrompt(query)

	response, err := s.provider.Chat([]interfaces.Message{
		{Role: "user", Content: prompt},
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate Cypher query: %w", err)
	}

	// Clean up response - extract just the Cypher query
	cypherQuery := strings.TrimSpace(response.Content)

	// Remove markdown code blocks if present
	if strings.HasPrefix(cypherQuery, "```") {
		// Extract content between code blocks
		lines := strings.Split(cypherQuery, "\n")
		var contentLines []string

		inCodeBlock := false
		for _, line := range lines {
			if strings.HasPrefix(line, "```") {
				inCodeBlock = !inCodeBlock
				continue // Skip the markers
			}

			if inCodeBlock || (!inCodeBlock && len(contentLines) > 0) {
				contentLines = append(contentLines, line)
			}
		}

		cypherQuery = strings.TrimSpace(strings.Join(contentLines, "\n"))
	}

	// Remove language identifier if present
	fmt.Println("Generated Cypher Query:", cypherQuery)
	cypherQuery = strings.TrimPrefix(cypherQuery, "cypher")
	cypherQuery = strings.TrimSpace(cypherQuery)

	return cypherQuery, nil
}

// TransformQueryToTopicCypher transforms a natural language topic query to a Cypher query that returns project IDs
func (s *DefaultQueryService) TransformQueryToTopicCypher(query string) (string, error) {
	prompt := utils.GetTopicSearchPrompt(query, s.graphRepo)
	fmt.Println("In TransformQueryToTopicCypher,Topic search query:", query)

	response, err := s.provider.Chat([]interfaces.Message{
		{Role: "user", Content: prompt},
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate topic search query: %w", err)
	}

	// Clean up response
	cypherQuery := strings.TrimSpace(response.Content)

	// Remove markdown code blocks if present
	if strings.Contains(cypherQuery, "```") {
		// Extract content between code blocks
		lines := strings.Split(cypherQuery, "\n")
		var contentLines []string

		inCodeBlock := false
		for _, line := range lines {
			if strings.HasPrefix(line, "```") {
				inCodeBlock = !inCodeBlock
				continue // Skip the markers
			}

			if inCodeBlock || (!inCodeBlock && len(contentLines) > 0) {
				contentLines = append(contentLines, line)
			}
		}

		cypherQuery = strings.TrimSpace(strings.Join(contentLines, "\n"))
	}

	return cypherQuery, nil
}

// ExecuteSearchQuery executes a search query and returns the results
func (s *DefaultQueryService) ExecuteSearchQuery(query string) ([]map[string]interface{}, error) {
	cypherQuery, err := s.TransformQueryToCypher(query)
	if err != nil {
		return nil, err
	}

	// Execute query using the graph repository
	results, err := s.graphRepo.ExecuteQuery(cypherQuery, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return results, nil
}

// ExecuteTopicSearchQuery executes a topic search query with focus on tags and related concepts
func (s *DefaultQueryService) ExecuteTopicSearchQuery(query string) ([]map[string]interface{}, error) {
	// Use the topic-specific transform function
	cypherQuery, err := s.TransformQueryToTopicCypher(query)
	if err != nil {
		return nil, err
	}

	fmt.Println("Executing topic search query:", cypherQuery)

	// Execute query using the graph repository
	results, err := s.graphRepo.ExecuteQuery(cypherQuery, map[string]interface{}{})
	if err != nil {
		fmt.Printf("Error executing topic query: %v\n", err)

		// Create a simplified tag-focused fallback query if the AI-generated one fails
		terms := strings.Split(query, " ")
		whereClauses := make([]string, 0, len(terms)*2)

		for _, term := range terms {
			if len(term) > 3 { // Skip short words
				whereClauses = append(whereClauses,
					fmt.Sprintf("toLower(t.name) CONTAINS toLower(\"%s\")", term))
			}
		}

		fallbackQuery := fmt.Sprintf(`
            MATCH (p:Project)-[:HAS_TAG]->(t:Tag)
            WHERE %s
            RETURN DISTINCT p.id as projectId, p.name as projectName
        `, strings.Join(whereClauses, " OR "))

		results, err = s.graphRepo.ExecuteQuery(fallbackQuery, map[string]interface{}{})
		if err != nil {
			return nil, fmt.Errorf("failed to execute fallback query: %w", err)
		}
	}

	return results, nil
}

// GetAllProjects fetches all projects with their basic information
func (s *DefaultQueryService) GetAllProjects() ([]map[string]interface{}, error) {
	// Simple Cypher query to fetch all projects with basic information
	// cypherQuery := `
	//     MATCH (p:Project)
	//     OPTIONAL MATCH (p)-[:HAS_TAG]->(t:Tag)
	//     RETURN p.id as projectId, p.name as projectName,
	//            p.description as description, p.status as status,
	//            collect(distinct t.name) as tags
	//     LIMIT 100
	// `
	cypherQuery := `MATCH (p:Project)
        OPTIONAL MATCH (p)-[:HAS_TAG]->(t:Tag)
        RETURN p.project_id as projectId, p.name as projectName, 
               p.description as description, 
               COALESCE(p.status, "published") as status,
               collect(distinct t.name) as tags
        LIMIT 100
	`

	fmt.Println("Fetching all projects")

	// Execute the query using the graph repository
	results, err := s.graphRepo.ExecuteQuery(cypherQuery, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all projects: %w", err)
	}
	fmt.Println("Fetched all projects:", results)

	return results, nil
}
