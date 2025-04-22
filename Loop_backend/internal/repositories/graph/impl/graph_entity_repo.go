package impl

import (
    "fmt"
    "Loop_backend/internal/models"
    "Loop_backend/platform/database/neo4j/entities"

    "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type entityRepo struct {
    baseRepo
}

// Identity-based vs Reusable node types
var identityNodes = map[string]bool{
    entities.TypeProject: true,
    entities.TypePerson:  true,
}

// GetProjectEntities retrieves all entities of a specific type for a project
func (r *entityRepo) GetProjectEntities(projectID string, entityType string) ([]models.Entity, error) {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := fmt.Sprintf(r.queryManager.GetEntityQuery("get_by_project_type"), entityType)
        records, err := tx.Run(query, map[string]interface{}{
            "project_id":  projectID,
        })
        if err != nil {
            return nil, err
        }

        var entities []models.Entity
        for records.Next() {
            record := records.Record()
            entities = append(entities, models.Entity{
                Name:        record.GetByIndex(0).(string),
                Type:        entityType,
                Description: record.GetByIndex(1).(string),
            })
        }
        return entities, nil
    })

    if err != nil {
        return nil, err
    }

    return result.([]models.Entity), nil
}

// GetEntitiesByType retrieves all entities of a specific type
func (r *entityRepo) GetEntitiesByType(entityType string) ([]models.Entity, error) {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := fmt.Sprintf(r.queryManager.GetEntityQuery("get_by_type"), entityType)
        records, err := tx.Run(query, map[string]interface{}{})
        if err != nil {
            return nil, err
        }

        var entities []models.Entity
        for records.Next() {
            record := records.Record()
            entities = append(entities, models.Entity{
                Name:        record.GetByIndex(0).(string),
                Type:        entityType,
                Description: record.GetByIndex(1).(string),
            })
        }
        return entities, nil
    })

    if err != nil {
        return nil, err
    }

    return result.([]models.Entity), nil
}

// CreateEntity creates a new entity in the graph
func (r *entityRepo) CreateEntity(entity *models.Entity, projectID string) error {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := r.queryManager.FormatEntityQuery(
            entity.Type,
            identityNodes[entity.Type],
        )
        params := map[string]interface{}{
            "project_id":  projectID,
            "name":        entity.Name,
            "description": entity.Description,
        }
        result, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume()
    })
    return err
}

// UpdateEntity updates an existing entity in the graph
func (r *entityRepo) UpdateEntity(entity *models.Entity) error {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := r.queryManager.FormatEntityQuery(entity.Type, false)
        params := map[string]interface{}{
            "name":        entity.Name,
            "description": entity.Description,
        }
        result, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume()
    })
    return err
}

// DeleteEntity deletes an entity from the graph
func (r *entityRepo) DeleteEntity(entityType, entityID string) error {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := r.queryManager.GetEntityQuery("delete")
        params := map[string]interface{}{
            "entity_type": entityType,
            "entity_id":   entityID,
        }
        result, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume()
    })
    return err
}
