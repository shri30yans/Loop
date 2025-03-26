package services

import (
    "Loop_backend/internal/ai/providers"
    "Loop_backend/internal/models"
    "Loop_backend/internal/utils"
    "fmt"
)

// DefaultProjectProcessor implements models.ProjectProcessor
type DefaultProjectProcessor struct {
    provider     providers.Provider
    tagService   TagService
    graphService GraphService
}

// NewProjectProcessor creates a new project analyzer instance
func NewProjectProcessor(
    provider providers.Provider,
    graphService GraphService,
    tagService TagService,
) models.ProjectProcessor {
    return &DefaultProjectProcessor{
        provider:     provider,
        graphService: graphService,
        tagService:   tagService,
    }
}

func (pa *DefaultProjectProcessor) fetchEntitiesAndRelations(project *models.Project) (*providers.ChatResponse, error) {
    prompt := utils.GetProjectAnalysisPrompt(project)
    
    // Convert project text to Entity Relationship
    response, err := pa.provider.Chat([]providers.Message{
        {Role: "user", Content: prompt},
    })
    
    if err != nil {
        fmt.Printf("Error sending prompt to LLM for project ID %s: %v\n", project.ProjectID, err)
        return nil, fmt.Errorf("failed to get LLM response: %w", err)
    }

    return response, nil
}

// AnalyzeNewProject performs LLM analysis on a new project
func (pa *DefaultProjectProcessor) AnalyzeNewProject(project *models.Project) error {
   response, err := pa.fetchEntitiesAndRelations(project)
    if err != nil {
        return err
    }
    
    // Parse entities and relations into knowledge graph 
    parser := utils.GetResponseParser()
    graph, err := parser.GenerateKnowledgeGraph(response.Content)
    if err != nil {
        fmt.Printf("Error parsing LLM response for project ID %s: %v\n", project.ProjectID, err)
        return fmt.Errorf("failed to parse response: %w", err)
    }
    fmt.Printf("Parsed knowledge graph for project ID %s: %+v\n", project.ProjectID, graph)

    // // Create tags from entities
    // for _, entity := range graph.Entities {
    //     tag := models.Tag{
    //         ProjectID:   project.ProjectID,
    //         Name:        entity.Name,
    //         Type:        entity.Type,
    //         Description: entity.Description,
    //         Category:    "entity",
    //     }
    //     err := pa.tagService.CreateTagWithEmbedding(&tag)
    //     if err != nil {
    //         fmt.Printf("Error creating entity tag for project ID %s: %+v, error: %v\n", project.ProjectID, tag, err)
    //         return fmt.Errorf("failed to create entity tag: %w", err)
    //     }
    // }

    // // Create relationship tags
    // for _, rel := range graph.Relationships {
    //     tag := models.Tag{
    //         ProjectID:   project.ProjectID,
    //         Name:        fmt.Sprintf("%s-%s", rel.Source, rel.Target),
    //         Type:        rel.Type,
    //         Description: rel.Description,
    //         Category:    rel.Category,
    //         UsageCount:  rel.Weight,
    //     }
    //     err := pa.tagService.CreateTagWithEmbedding(&tag)
    //     if err != nil {
    //         fmt.Printf("Error creating relationship tag for project ID %s: %+v, error: %v\n", project.ProjectID, tag, err)
    //         return fmt.Errorf("failed to create relationship tag: %w", err)
    //     }
    // }

    // Store relationships in graph database
    err = pa.graphService.StoreProjectGraph(project, graph)
    if err != nil {
        fmt.Printf("Error storing relationships in graph database for project ID %s: %v\n", project.ProjectID, err)
        return fmt.Errorf("failed to store relationships: %w", err)
    }
    fmt.Printf("Successfully stored relationships in graph database for project ID %s\n", project.ProjectID)

    // Log the successful completion of the analysis
    fmt.Printf("Completed analysis for project ID: %s\n", project.ProjectID)

    return nil
}
