package services

import (
    "Loop_backend/internal/models"
    "strings"
)

type TagGenerationService interface {
    GenerateProjectTags(project *models.Project) ([]string, error)
    MergeTags(userTags []string, generatedTags []string) []string
}

type tagGenerationService struct {
    textProcessor TextProcessor
}

func NewTagGenerationService(textProcessor TextProcessor) TagGenerationService {
    return &tagGenerationService{
        textProcessor: textProcessor,
    }
}

func (s *tagGenerationService) GenerateProjectTags(project *models.Project) ([]string, error) {
    // Combine project text fields
    var texts []string
    texts = append(texts, project.Title, project.Description, project.Introduction)
    for _, section := range project.Sections {
        texts = append(texts, section.Title, section.Body)
    }
    combinedText := strings.Join(texts, " ")

    // Process text to extract tags
    tags := s.textProcessor.ExtractKeywords(combinedText)
    return tags, nil
}

// MergeTags combines user-provided tags with generated tags, removing duplicates
func (s *tagGenerationService) MergeTags(userTags []string, generatedTags []string) []string {
    // Create a map to track unique tags
    tagSet := make(map[string]bool)
    
    // Add all user tags first (they take priority)
    for _, tag := range userTags {
        tagSet[tag] = true
    }
    
    // Add generated tags if they don't already exist
    for _, tag := range generatedTags {
        tagSet[tag] = true
    }
    
    // Convert map keys back to slice
    mergedTags := make([]string, 0, len(tagSet))
    for tag := range tagSet {
        mergedTags = append(mergedTags, tag)
    }
    
    return mergedTags
}
