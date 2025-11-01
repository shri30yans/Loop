package project

import (
	"Loop_backend/internal/models"
)

// ProjectProcessor defines the interface for project analysis
type ProjectProcessor interface {
	AnalyzeNewProject(project *models.Project) error
	GenerateProjectSummary(project *models.Project) (string, error)
}