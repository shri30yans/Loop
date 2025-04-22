package impl

import (
	"Loop_backend/internal/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type projectRepo struct {
	baseRepo
}

func debugPrintParams(operation string, params map[string]interface{}) {
	fmt.Printf("\n=== Neo4j Operation: %s ===\n", operation)
	fmt.Printf("Parameters: %+v\n", params)
	fmt.Println("======================")
}

// StoreProjectGraph stores the entire knowledge graph for a project
func (r *projectRepo) StoreProjectGraph(projectID uuid.UUID, ownerID string, graph *models.KnowledgeGraph) error {
	fmt.Printf("\n\n=== Storing Project Graph: %s ===\n", projectID)
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Create project node and link owner
		projectQuery := r.queryManager.GetProjectQuery("create")
		fmt.Printf("\nExecuting Create Project Query:\n%s\n", projectQuery)

		// Find the project entity in the graph
		var projectName, projectDesc string
		for _, entity := range graph.Entities {
			if entity.Type == "Project" {
				projectName = entity.Name
				projectDesc = entity.Description
				break
			}
		}

		params := map[string]interface{}{
			"project_id":  projectID.String(),
			"name":        projectName,
			"description": projectDesc,
		}
		debugPrintParams("Create Project", params)
		if _, err := tx.Run(projectQuery, params); err != nil {
			fmt.Printf("Error creating project node: %v\n", err)
			return nil, err
		}

		ownerQuery := r.queryManager.GetProjectQuery("link_owner")
		fmt.Printf("\nExecuting Link Owner Query:\n%s\n", ownerQuery)
		params = map[string]interface{}{
			"owner_id":   ownerID,
			"project_id": projectID.String(),
		}
		debugPrintParams("Link Owner", params)
		if _, err := tx.Run(ownerQuery, params); err != nil {
			fmt.Printf("Error linking owner: %v\n", err)
			return nil, err
		}

		// Group entities by type for batch processing
		entityTypeMap := make(map[string][]map[string]interface{})
		for _, entity := range graph.Entities {
			entityTypeMap[entity.Type] = append(entityTypeMap[entity.Type], map[string]interface{}{
				"type":        entity.Type,
				"name":        entity.Name,
				"description": entity.Description,
				"identity":    identityNodes[entity.Type],
			})
		}

		// Process each type's entities in a batch
		for entityType, entities := range entityTypeMap {
			fmt.Printf("\n=== Batch storing %d %s entities ===\n", len(entities), entityType)
			batchQuery := fmt.Sprintf(r.queryManager.GetEntityQuery("batch_create"), entityType)
			batchParams := map[string]interface{}{
				"project_id": projectID.String(),
				"entities":   entities,
			}
			debugPrintParams(fmt.Sprintf("Batch Create %s Entities", entityType), batchParams)
			if _, err := tx.Run(batchQuery, batchParams); err != nil {
				fmt.Printf("Error batch creating %s entities: %v\n", entityType, err)
				return nil, err
			}
		}

		// Store relationships
		fmt.Printf("\n=== Storing %d Relationships ===\n", len(graph.Relationships))
		for i, rel := range graph.Relationships {
			var sourceType, targetType string
			for _, entity := range graph.Entities {
				if entity.Name == rel.Source {
					sourceType = entity.Type
				}
				if entity.Name == rel.Target {
					targetType = entity.Type
				}
				if sourceType != "" && targetType != "" {
					break
				}
			}

			useIdentityBased := identityNodes[sourceType] && identityNodes[targetType]
			query := r.queryManager.FormatCreateRelationship(
				sourceType,
				targetType,
				rel.Type,
				useIdentityBased,
			)
			fmt.Printf("\nExecuting Relationship Query %d:\n%s\n", i+1, query)

			params := map[string]interface{}{
				"source":      rel.Source,
				"target":      rel.Target,
				"source_id":   projectID.String(),
				"target_id":   projectID.String(),
				"description": rel.Description,
				"weight":      rel.Weight,
				"category":    rel.Category,
			}
			debugPrintParams(fmt.Sprintf("Create Relationship %d (%s -> %s)", i+1, rel.Source, rel.Target), params)
			if _, err := tx.Run(query, params); err != nil {
				fmt.Printf("Error creating relationship: %v\n", err)
				return nil, err
			}
		}

		return nil, nil
	})

	return err
}

// GetProjectGraph retrieves a project's knowledge graph
func (r *projectRepo) GetProjectGraph(projectID uuid.UUID) (*models.KnowledgeGraph, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Get nodes
		nodesQuery := r.queryManager.GetProjectQuery("get_nodes")
		entitiesResult, err := tx.Run(nodesQuery, map[string]interface{}{
			"project_id": projectID.String(),
		})
		if err != nil {
			return nil, err
		}

		var entities []models.Entity
		for entitiesResult.Next() {
			record := entitiesResult.Record()
			entities = append(entities, models.Entity{
				Name:        record.GetByIndex(0).(string),
				Type:        record.GetByIndex(1).(string),
				Description: record.GetByIndex(2).(string),
			})
		}

		// Get relationships
		relsQuery := r.queryManager.GetProjectQuery("get_relationships")
		relationshipsResult, err := tx.Run(relsQuery, map[string]interface{}{
			"project_id": projectID.String(),
		})
		if err != nil {
			return nil, err
		}

		var relationships []models.Relationship
		for relationshipsResult.Next() {
			record := relationshipsResult.Record()
			weight := 5 // default weight
			if w, ok := record.GetByIndex(4).(int64); ok {
				weight = int(w)
			}

			relationships = append(relationships, models.Relationship{
				Source:      record.GetByIndex(0).(string),
				Target:      record.GetByIndex(1).(string),
				Type:        record.GetByIndex(2).(string),
				Description: record.GetByIndex(3).(string),
				Weight:      weight,
				Category:    record.GetByIndex(5).(string),
			})
		}

		return &models.KnowledgeGraph{
			Entities:      entities,
			Relationships: relationships,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.KnowledgeGraph), nil
}

// DeleteProjectNode deletes a project and all its relationships
func (r *projectRepo) DeleteProjectNode(projectID string) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := r.queryManager.GetProjectQuery("delete")
		result, err := tx.Run(query, map[string]interface{}{
			"project_id": projectID,
		})
		if err != nil {
			return nil, err
		}
		return result.Consume()
	})
	return err
}
