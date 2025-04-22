package services

import (
	"Loop_backend/internal/ai/interfaces"
	"Loop_backend/internal/repositories"
	"Loop_backend/internal/utils"
	"fmt"

	"github.com/google/uuid"
)

// SummaryService defines operations for generating project summaries
type SummaryService interface {
	GenerateProjectSummary(projectID string) (map[string]interface{}, error)
}

// DefaultSummaryService implements SummaryService
type DefaultSummaryService struct {
	provider    interfaces.Provider
	projectRepo repositories.ProjectRepository
	graphRepo   repositories.GraphRepository
}

// NewSummaryService creates a new summary service
func NewSummaryService(
	provider interfaces.Provider,
	projectRepo repositories.ProjectRepository,
	graphRepo repositories.GraphRepository,
) SummaryService {
	return &DefaultSummaryService{
		provider:    provider,
		projectRepo: projectRepo,
		graphRepo:   graphRepo,
	}
}

// GenerateProjectSummary fetches project data and generates a summary
func (s *DefaultSummaryService) GenerateProjectSummary(projectID string) (map[string]interface{}, error) {
	// Parse project ID
	id, err := uuid.Parse(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID: %w", err)
	}

	// Fetch project from SQL database using the correct method
	project, err := s.projectRepo.GetProject(id.String())
	if err != nil {
		return nil, fmt.Errorf("error fetching project: %w", err)
	}

	// Remove the problematic code and just rely on your direct query
	var projectDetails map[string]interface{}
	projectDetails = map[string]interface{}{
		"name": project.Title,
	}
	fmt.Println("Project Details:", projectDetails)

	// Use parameterized query
	relatedEntitiesQuery := `
        MATCH (p:Project)
        WHERE p.id = $projectId
        OPTIONAL MATCH (p)-[:USES]->(tech:Technology)
        OPTIONAL MATCH (p)-[:HAS_STAKEHOLDER]->(stake:Stakeholder)
        OPTIONAL MATCH (p)-[:BELONGS_TO]->(cat:Category)
        RETURN collect(DISTINCT tech.name) AS technologies, 
               collect(DISTINCT stake.name) AS stakeholders,
               collect(DISTINCT cat.name) AS categories
    `

	// Use parameters map
	params := map[string]interface{}{
		"projectId": projectID,
	}

	technologies := make([]string, 0)
	stakeholders := make([]string, 0)
	categories := make([]string, 0)

	results, err := s.graphRepo.ExecuteQuery(relatedEntitiesQuery, params)
	if err == nil && len(results) > 0 {
		// Process technologies
		if techList, ok := results[0]["technologies"].([]interface{}); ok {
			for _, tech := range techList {
				if techName, ok := tech.(string); ok {
					technologies = append(technologies, techName)
				}
			}
		}

		// Process stakeholders
		if stakeList, ok := results[0]["stakeholders"].([]interface{}); ok {
			for _, stake := range stakeList {
				if stakeName, ok := stake.(string); ok {
					stakeholders = append(stakeholders, stakeName)
				}
			}
		}

		// Process categories
		if catList, ok := results[0]["categories"].([]interface{}); ok {
			for _, cat := range catList {
				if catName, ok := cat.(string); ok {
					categories = append(categories, catName)
				}
			}
		}
	}

	// Generate prompt with project information
	prompt := utils.GetProjectSummaryPrompt(project, technologies, stakeholders)

	// Call AI provider for summary
	response, err := s.provider.Chat([]interfaces.Message{
		{Role: "user", Content: prompt},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}

	// Return result
	return map[string]interface{}{
		"project":      project,
		"technologies": technologies,
		"stakeholders": stakeholders,
		"categories":   categories,
		"summary":      response.Content,
	}, nil
}
