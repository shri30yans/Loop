package project

import (
    "fmt"
    "strings"
    "Loop_backend/internal/models"
    "Loop_backend/internal/repositories/vector"
    "Loop_backend/internal/services/graph"
    "Loop_backend/internal/services/embedding"
    "Loop_backend/internal/ai/providers"
    "Loop_backend/internal/services/graph/response_parser"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

// DefaultProjectProcessor implements ProjectProcessor
type DefaultProjectProcessor struct {
    graphService    graph.Service
    vectorRepo      vector.VectorRepository
    provider        providers.Provider
    embeddingClient embedding.EmbeddingClient
    responseParser  response_parser.ResponseParser
    logger          zerolog.Logger
}

// NewProjectProcessor creates a new project processor
func NewProjectProcessor(
    graphService graph.Service,
    vectorRepo vector.VectorRepository,
    provider providers.Provider,
    embeddingClient embedding.EmbeddingClient,
    responseParser response_parser.ResponseParser,
) *DefaultProjectProcessor {
    return &DefaultProjectProcessor{
        graphService:    graphService,
        vectorRepo:      vectorRepo,
        provider:        provider,
        embeddingClient: embeddingClient,
        responseParser:  responseParser,
        logger:          log.With().Str("component", "ProjectProcessor").Logger(),
    }
}

// AnalyzeNewProject performs AI analysis on a new project
func (pa *DefaultProjectProcessor) AnalyzeNewProject(project *models.Project) error {
    pa.logger.Info().Str("projectID", project.ProjectID.String()).Msg("Processing project")

    // Extract entities using AI
    response, err := pa.extractEntities(project)
    if err != nil {
        pa.logger.Error().Err(err).Str("projectID", project.ProjectID.String()).Msg("Failed to extract entities")
        return err
    }

    // First extract and clean the summary
    summary := pa.extractSummary(response.Content)
    if summary == "" {
        pa.logger.Error().Str("projectID", project.ProjectID.String()).Msg("Failed to extract summary from response")
        return fmt.Errorf("failed to extract summary from response")
    }

    // Clean the response by removing the summary section
    cleanedContent := pa.cleanResponse(response.Content)

    // Parse cleaned response into knowledge graph
    graph, err := pa.responseParser.GenerateKnowledgeGraph(cleanedContent)
    if err != nil {
        pa.logger.Error().Err(err).Str("projectID", project.ProjectID.String()).Msg("Failed to parse response")
        return fmt.Errorf("failed to parse response: %w", err)
    }

    // Generate vector for summary
    vector, err := pa.embeddingClient.GenerateEmbedding(summary)
    if err != nil {
        pa.logger.Error().Err(err).Str("projectID", project.ProjectID.String()).Msg("Failed to generate embedding")
        return fmt.Errorf("failed to generate embedding: %w", err)
    }

    // Store knowledge graph
    err = pa.graphService.StoreProjectGraph(project.ProjectID.String(), graph)
    if err != nil {
        pa.logger.Error().Err(err).Str("projectID", project.ProjectID.String()).Msg("Failed to store graph")
        return fmt.Errorf("failed to store graph: %w", err)
    }

    // Store vector and summary
    err = pa.vectorRepo.StoreProjectVector(project.ProjectID, summary, vector)
    if err != nil {
        pa.logger.Error().Err(err).Str("projectID", project.ProjectID.String()).Msg("Failed to store vector")
        return fmt.Errorf("failed to store vector: %w", err)
    }
    
    pa.logger.Info().Str("projectID", project.ProjectID.String()).Msg("Project processing completed successfully")

    return nil
}

// extractEntities handles the AI-based entity extraction
func (pa *DefaultProjectProcessor) extractEntities(project *models.Project) (*providers.ChatResponse, error) {
    prompt := GetProjectAnalysisPrompt(project)
    pa.logger.Debug().Str("projectID", project.ProjectID.String()).Msg("Extracting entities with AI")
    response, err := pa.provider.Chat([]providers.Message{
        {Role: "user", Content: prompt},
    })
    if err != nil {
        pa.logger.Error().Err(err).Str("projectID", project.ProjectID.String()).Msg("Failed to get AI response")
        return nil, fmt.Errorf("failed to get AI response: %w", err)
    }
    return response, nil
}

// extractSummary extracts the summary from between summary tags
func (pa *DefaultProjectProcessor) extractSummary(content string) string {
    startTag := "[SUMMARY_START]"
    endTag := "[SUMMARY_END]"
    
    startIndex := strings.Index(content, startTag)
    if startIndex == -1 {
        return ""
    }
    startIndex += len(startTag)
    
    endIndex := strings.Index(content[startIndex:], endTag)
    if endIndex == -1 {
        return ""
    }
    
    return strings.TrimSpace(content[startIndex : startIndex+endIndex])
}

// cleanResponse removes the summary section from the content
func (pa *DefaultProjectProcessor) cleanResponse(content string) string {
    startTag := "[SUMMARY_START]"
    endTag := "[SUMMARY_END]"
    
    startIndex := strings.Index(content, startTag)
    if startIndex == -1 {
        return content
    }
    
    endIndex := strings.Index(content[startIndex:], endTag)
    if endIndex == -1 {
        return content
    }
    endIndex += startIndex + len(endTag)
    
    // Remove the summary section and any surrounding whitespace
    return strings.TrimSpace(content[:startIndex] + content[endIndex:])
}

// GenerateProjectSummary generates a summary for the project
func (pa *DefaultProjectProcessor) GenerateProjectSummary(project *models.Project) (string, error) {
    response, err := pa.extractEntities(project)
    if err != nil {
        return "", err
    }
    
    return pa.extractSummary(response.Content), nil
}
