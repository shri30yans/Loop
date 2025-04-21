package utils

import (
	"Loop_backend/internal/models"
	"fmt"
	"strings"
)

// GetProjectSummaryPrompt generates a prompt for summarizing project details
func GetProjectSummaryPrompt(project *models.Project, technologies []string, stakeholders []string) string {
	techStack := "None specified"
	if len(technologies) > 0 {
		techStack = strings.Join(technologies, ", ")
	}

	stakeholderList := "None specified"
	if len(stakeholders) > 0 {
		stakeholderList = strings.Join(stakeholders, ", ")
	}

	return fmt.Sprintf(`
---Project Summary Task---
Generate a comprehensive analysis and summary of this project based on the provided details.

---Project Information---
Name: %s
Description: %s
Introduction: %s
Status: %s
Tags: %s

---Technical Stack---
Technologies: %s

---Key Stakeholders---
Stakeholders: %s

---Analysis Requirements---
1. Provide a brief one-paragraph overview of the project
2. Analyze the technical stack and its appropriateness for the project goals
3. Identify potential challenges or opportunities based on the project description
4. Suggest possible enhancements or next steps
5. Estimate the project's potential impact in its domain

Format your response as a structured report with clear sections.
`,
		project.Title,
		project.Description,
		project.Introduction,
		string(project.Status),
		strings.Join(project.Tags, ", "),
		techStack,
		stakeholderList,
	)
}
