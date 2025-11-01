package search

import (
	"Loop_backend/internal/ai/interfaces"
	"Loop_backend/internal/dto"
	graphRepo "Loop_backend/internal/repositories/graph"
	"fmt"
	"strings"
)

// CypherSearchService defines operations for Cypher-based search
type CypherSearchService interface {
	Search(query string) ([]CypherSearchToolsResult, error)
	SearchProjects(query string) ([]dto.SearchToolsResult, error)
}

// CypherSearchToolsResult represents a search result from Neo4j
type CypherSearchToolsResult struct {
	ProjectID   string `json:"projectId"`
	ProjectName string `json:"projectName"`
	Summary     string `json:"summary"`
}

// DefaultCypherSearchService implements CypherSearchService
type DefaultCypherSearchService struct {
	provider  interfaces.Provider
	graphRepo graphRepo.Repository
}

// NewCypherSearchService creates a new cypher search service instance
func NewCypherSearchService(
	provider interfaces.Provider,
	graphRepo graphRepo.Repository,
) CypherSearchService {
	return &DefaultCypherSearchService{
		provider:  provider,
		graphRepo: graphRepo,
	}
}

// SearchProjects performs a search and returns project search results
func (s *DefaultCypherSearchService) SearchProjects(query string) ([]dto.SearchToolsResult, error) {
	results, err := s.Search(query)
	if err != nil {
		return nil, err
	}

	SearchToolsResults := make([]dto.SearchToolsResult, 0, len(results))
	for _, result := range results {
		SearchToolsResults = append(SearchToolsResults, dto.SearchToolsResult{
			ProjectID:   result.ProjectID,
			Name:        result.ProjectName,
			Description: result.Summary,
		})
	}

	return SearchToolsResults, nil
}

// Search performs a search by converting the query to Cypher and executing it
func (s *DefaultCypherSearchService) Search(query string) ([]CypherSearchToolsResult, error) {
	// Debug: Print input query
	fmt.Printf("\nProcessing Search Query: %s\n", query)

	// Transform query to Cypher
	cypherQuery, err := s.transformQueryToCypher(query)
	if err != nil {
		return nil, fmt.Errorf("failed to transform query: %w", err)
	}

	// Debug: Print generated Cypher query
	fmt.Printf("Generated Cypher Query: %s\n", cypherQuery)

	// Execute the Cypher query
	results, err := s.graphRepo.ExecuteQuery(cypherQuery, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// Convert results to CypherSearchToolsResult
	SearchToolsResults := make([]CypherSearchToolsResult, 0, len(results))
	for _, result := range results {
		projectID, ok := result["projectId"].(string)
		if !ok {
			continue
		}

		projectName, ok := result["projectName"].(string)
		if !ok {
			projectName = "" // Use empty string if name not found
		}

		// Debug: Print found project ID
		fmt.Printf("Found Project - ID: %s, Name: %s\n", projectID, projectName)

		// Use project description as summary
		summary := ""
		if desc, ok := result["description"].(string); ok {
			summary = desc
		}

		SearchToolsResults = append(SearchToolsResults, CypherSearchToolsResult{
			ProjectID:   projectID,
			ProjectName: projectName,
			Summary:     summary,
		})
	}

	// Debug: Print total results found
	fmt.Printf("Total Results Found: %d\n\n", len(SearchToolsResults))

	return SearchToolsResults, nil
}

// transformQueryToCypher transforms a natural language query to a Cypher query
func (s *DefaultCypherSearchService) transformQueryToCypher(query string) (string, error) {
	prompt := GetCypherTransformPrompt(query)

	// Debug: Print the prompt being used
	fmt.Printf("Using Transformation Prompt:\n%s\n", prompt)

	response, err := s.provider.Chat([]interfaces.Message{
		{Role: "user", Content: prompt},
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate Cypher query: %w", err)
	}

	return cleanCypherQuery(response.Content), nil
}

// cleanCypherQuery cleans up the Cypher query from AI response
func cleanCypherQuery(query string) string {
	query = strings.TrimSpace(query)

	// Remove markdown code blocks if present
	if strings.HasPrefix(query, "```") {
		// Extract content between code blocks
		lines := strings.Split(query, "\n")
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

		query = strings.TrimSpace(strings.Join(contentLines, "\n"))
	}

	// Remove language identifier if present
	query = strings.TrimPrefix(query, "cypher")
	query = strings.TrimSpace(query)

	return query
}
