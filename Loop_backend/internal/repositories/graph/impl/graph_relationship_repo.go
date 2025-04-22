package impl

import (
    "Loop_backend/internal/models"

    "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type relationshipRepo struct {
    baseRepo
}

// CreateRelationship creates a relationship between two entities
func (r *relationshipRepo) CreateRelationship(source, target string, relType string, props map[string]interface{}) error {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := r.queryManager.GetRelationshipQuery("create")
        params := map[string]interface{}{
            "source": source,
            "target": target,
            "type":   relType,
            "props":  props,
        }
        result, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume()
    })
    return err
}

// GetProjectRelationships retrieves all relationships for a project
func (r *relationshipRepo) GetProjectRelationships(projectID string) ([]models.Relationship, error) {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := r.queryManager.GetProjectQuery("get_relationships")
        records, err := tx.Run(query, map[string]interface{}{
            "project_id": projectID,
        })
        if err != nil {
            return nil, err
        }

        var relationships []models.Relationship
        for records.Next() {
            record := records.Record()
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
        return relationships, nil
    })

    if err != nil {
        return nil, err
    }

    return result.([]models.Relationship), nil
}

// GetEntityRelationships retrieves all relationships for an entity
func (r *relationshipRepo) GetEntityRelationships(entityType, entityID string) ([]models.Relationship, error) {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := r.queryManager.GetRelationshipQuery("get_all")
        records, err := tx.Run(query, map[string]interface{}{
            "entity_id": entityID,
        })
        if err != nil {
            return nil, err
        }

        var relationships []models.Relationship
        for records.Next() {
            record := records.Record()
            props := record.GetByIndex(3).(map[string]interface{})
            weight := 5
            if w, ok := props["weight"].(int64); ok {
                weight = int(w)
            }

            relationships = append(relationships, models.Relationship{
                Type:        record.GetByIndex(0).(string),
                Source:      entityID,
                Target:      record.GetByIndex(2).(string),
                Description: props["description"].(string),
                Weight:      weight,
                Category:    props["category"].(string),
            })
        }
        return relationships, nil
    })

    if err != nil {
        return nil, err
    }

    return result.([]models.Relationship), nil
}

// UpdateRelationship updates the properties of a relationship
func (r *relationshipRepo) UpdateRelationship(relationshipID string, props map[string]interface{}) error {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := r.queryManager.GetRelationshipQuery("update")
        params := map[string]interface{}{
            "relationship_id": relationshipID,
            "props":          props,
        }
        result, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume()
    })
    return err
}

// DeleteRelationship deletes a relationship by ID
func (r *relationshipRepo) DeleteRelationship(relationshipID string) error {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := r.queryManager.GetRelationshipQuery("delete")
        params := map[string]interface{}{
            "relationship_id": relationshipID,
        }
        result, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume()
    })
    return err
}
