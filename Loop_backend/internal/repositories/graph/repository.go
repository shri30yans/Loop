package graph

import (
    "context"
    "fmt"
    "strings"
    "Loop_backend/internal/models"
    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Repository defines the interface for graph operations
type Repository interface {
    // Search operations
    GetProjectIDsByTags(tags []string) ([]string, error)
    GetProjectIDsByText(searchText string) ([]string, error)
    ExecuteQuery(query string, params map[string]interface{}) ([]map[string]interface{}, error)
    
    // Graph operations
    StoreProjectGraph(projectID string, graph *models.KnowledgeGraph) error
}

// Neo4jRepository implements Repository using Neo4j
type Neo4jRepository struct {
    driver neo4j.DriverWithContext
}

// NewRepository creates a new Neo4j repository
func NewRepository(driver neo4j.DriverWithContext) Repository {
    return &Neo4jRepository{
        driver: driver,
    }
}

// GetProjectIDsByTags returns project IDs that have any of the specified tags
func (r *Neo4jRepository) GetProjectIDsByTags(tags []string) ([]string, error) {
    ctx := context.Background()
    session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
    defer session.Close(ctx)

    result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
        query := `
            MATCH (p:Project)-[:HAS_TAG]->(t:Tag)
            WHERE toLower(t.name) IN $tags
            RETURN DISTINCT p.id as projectId`

        // Convert tags to lowercase for case-insensitive matching
        lowerTags := make([]string, len(tags))
        for i, tag := range tags {
            lowerTags[i] = strings.ToLower(tag)
        }

        result, err := tx.Run(ctx, query, map[string]interface{}{
            "tags": lowerTags,
        })
        if err != nil {
            return nil, err
        }

        var projectIDs []string
        records, err := result.Collect(ctx)
        if err != nil {
            return nil, err
        }

        for _, record := range records {
            if projectID, ok := record.Get("projectId"); ok {
                projectIDs = append(projectIDs, projectID.(string))
            }
        }

        return projectIDs, nil
    })

    if err != nil {
        return nil, fmt.Errorf("failed to get projects by tags: %w", err)
    }

    return result.([]string), nil
}

// GetProjectIDsByText returns project IDs based on text search
func (r *Neo4jRepository) GetProjectIDsByText(searchText string) ([]string, error) {
    ctx := context.Background()
    session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
    defer session.Close(ctx)

    result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
        query := `
            MATCH (p:Project)
            WHERE toLower(p.name) CONTAINS toLower($text)
            OR toLower(p.description) CONTAINS toLower($text)
            RETURN DISTINCT p.id as projectId`

        result, err := tx.Run(ctx, query, map[string]interface{}{
            "text": searchText,
        })
        if err != nil {
            return nil, err
        }

        var projectIDs []string
        records, err := result.Collect(ctx)
        if err != nil {
            return nil, err
        }

        for _, record := range records {
            if projectID, ok := record.Get("projectId"); ok {
                projectIDs = append(projectIDs, projectID.(string))
            }
        }

        return projectIDs, nil
    })

    if err != nil {
        return nil, fmt.Errorf("failed to get projects by text: %w", err)
    }

    return result.([]string), nil
}

// ExecuteQuery executes a custom Cypher query
func (r *Neo4jRepository) ExecuteQuery(query string, params map[string]interface{}) ([]map[string]interface{}, error) {
    ctx := context.Background()
    session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
    defer session.Close(ctx)

    result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
        result, err := tx.Run(ctx, query, params)
        if err != nil {
            return nil, err
        }

        records, err := result.Collect(ctx)
        if err != nil {
            return nil, err
        }

        var resultRecords []map[string]interface{}
        for _, record := range records {
            recordMap := make(map[string]interface{})
            for i, key := range record.Keys {
                value := record.Values[i]
                if node, ok := value.(neo4j.Node); ok {
                    recordMap[key] = node.Props
                } else {
                    recordMap[key] = value
                }
            }
            resultRecords = append(resultRecords, recordMap)
        }
        return resultRecords, nil
    })

    if err != nil {
        return nil, fmt.Errorf("failed to execute query: %w", err)
    }

    return result.([]map[string]interface{}), nil
}

// StoreProjectGraph stores the knowledge graph for a project
func (r *Neo4jRepository) StoreProjectGraph(projectID string, graph *models.KnowledgeGraph) error {
    ctx := context.Background()
    session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
    defer session.Close(ctx)

    _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
        // Create or update project node
        projectQuery := `
            MERGE (p:Project {id: $projectId})
            SET p.name = $name,
                p.description = $description,
                p.keywords = $keywords
        `
        _, err := tx.Run(ctx, projectQuery, map[string]interface{}{
            "projectId":   projectID,
            "name":        graph.Name,
            "description": graph.Description,
            "keywords":    graph.Keywords,
        })
        if err != nil {
            return nil, fmt.Errorf("failed to create project node: %w", err)
        }

        // Create or update entities
        for _, entity := range graph.Entities {
            entityQuery := `
                MERGE (e:Entity {name: $name, type: $type})
                SET e.level = $level,
                    e.description = $description
                WITH e
                MATCH (p:Project {id: $projectId})
                MERGE (p)-[:HAS_ENTITY]->(e)
            `
            _, err := tx.Run(ctx, entityQuery, map[string]interface{}{
                "projectId":   projectID,
                "name":        entity.Name,
                "type":        entity.Type,
                "level":       entity.Level,
                "description": entity.Description,
            })
            if err != nil {
                return nil, fmt.Errorf("failed to create entity node: %w", err)
            }

            // Handle parent-child relationships for entities
            if entity.ParentEntity != "" {
                parentQuery := `
                    MATCH (child:Entity {name: $childName})
                    MATCH (parent:Entity {name: $parentName})
                    MERGE (parent)-[:CONTAINS]->(child)
                `
                _, err := tx.Run(ctx, parentQuery, map[string]interface{}{
                    "childName":  entity.Name,
                    "parentName": entity.ParentEntity,
                })
                if err != nil {
                    return nil, fmt.Errorf("failed to create parent-child relationship: %w", err)
                }
            }
        }

        // Create relationships
        for _, rel := range graph.Relationships {
            relationshipQuery := `
                MATCH (source:Entity {name: $sourceName})
                MATCH (target:Entity {name: $targetName})
                MERGE (source)-[r:${relType} {weight: $weight}]->(target)
            `
            // Replace template with actual relationship type
            relationshipQuery = strings.Replace(relationshipQuery, "${relType}", rel.Type, 1)

            _, err := tx.Run(ctx, relationshipQuery, map[string]interface{}{
                "sourceName": rel.Source,
                "targetName": rel.Target,
                "weight":     rel.Weight,
            })
            if err != nil {
                return nil, fmt.Errorf("failed to create relationship: %w", err)
            }
        }

        return nil, nil
    })

    if err != nil {
        return fmt.Errorf("failed to store project graph: %w", err)
    }

    return nil
}
