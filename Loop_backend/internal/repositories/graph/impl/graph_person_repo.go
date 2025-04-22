package impl

import (
    "Loop_backend/internal/models"
    "time"

    "github.com/google/uuid"
    "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type personRepo struct {
    baseRepo
}

// GetPersonProjects retrieves all projects associated with a person
func (r *personRepo) GetPersonProjects(personID string) ([]models.Project, error) {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := r.queryManager.GetPersonQuery("get_projects")
        records, err := tx.Run(query, map[string]interface{}{
            "userID": personID,
        })
        if err != nil {
            return nil, err
        }

        var projects []models.Project
        for records.Next() {
            record := records.Record()
            projectID, _ := uuid.Parse(record.GetByIndex(0).(string))
            projects = append(projects, models.Project{
                ProjectInfo: models.ProjectInfo{
                    ProjectID:   projectID,
                    Title:       record.GetByIndex(1).(string),
                    Description: record.GetByIndex(2).(string),
                },
            })
        }
        return projects, nil
    })

    if err != nil {
        return nil, err
    }

    return result.([]models.Project), nil
}

// GetPersonContributions retrieves all contributions made by a person
func (r *personRepo) GetPersonContributions(personID string, limit int) ([]models.Contribution, error) {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := r.queryManager.GetPersonQuery("get_contributions")
        records, err := tx.Run(query, map[string]interface{}{
            "userID": personID,
            "limit":  limit,
        })
        if err != nil {
            return nil, err
        }

        var contributions []models.Contribution
        for records.Next() {
            record := records.Record()
            contributions = append(contributions, models.Contribution{
                ActivityType: record.GetByIndex(0).(string),
                TargetName:  record.GetByIndex(1).(string),
                TargetType:  record.GetByIndex(2).(string),
                Timestamp:   time.Unix(record.GetByIndex(3).(int64), 0),
            })
        }
        return contributions, nil
    })

    if err != nil {
        return nil, err
    }

    return result.([]models.Contribution), nil
}

// LinkPersonToEntity creates a relationship between a person and an entity
func (r *personRepo) LinkPersonToEntity(personID, entityID string, entityType string, relationType string, props map[string]interface{}) error {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := r.queryManager.FormatPersonEntityLink(entityType, "id", relationType)
        params := map[string]interface{}{
            "userID":     personID,
            "entityID":   entityID,
            "properties": props,
        }
        result, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume()
    })
    return err
}
