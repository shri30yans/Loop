package repositories

import (
    "time"
    "Loop_backend/internal/models"
    "Loop_backend/internal/ai/processor"
    "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type GraphRepository interface {
    // Existing project operations
    CreateProjectWithUserAndTags(project *models.Project, tags []string) error
    GetProjectTags(projectID string) ([]string, error)
    GetProjectsWithTags(projectIDs []string) (map[string][]string, error)
    DeleteProjectNode(projectID string) error

    // New tag operations
    CreateTagNode(tag *models.Tag) error
    UpdateTagNode(tag *models.Tag) error
    CreateTagRelationship(tag1, tag2 string, strength float64) error
    GetRelatedTags(tagName string, minStrength float64) ([]*models.TagRelationship, error)
    
    // User-Tag operations
    SetUserTagExpertise(userID, tagName string, level string, years int) error
    GetTagExperts(tagName string) ([]models.User, error)
    GetUserExpertise(userID string) (map[string]string, error)
    
    // Enhanced project queries
    GetProjectsByTag(tagName string) ([]*models.Project, error)
    GetUsersWithTag(tagName string) ([]string, error)
    GetProjectCollaborators(projectID string) ([]string, error)

    // Entity operations
    CreateEntitiesAndRelationships(projectID string, entities []processor.Entity, relationships []processor.Relationship) error
}

type graphRepository struct {
    driver neo4j.Driver
}

func NewGraphRepository(driver neo4j.Driver) GraphRepository {
    if driver == nil {
        return nil
    }
    return &graphRepository{
        driver: driver,
    }
}

func (r *graphRepository) CreateProjectWithUserAndTags(project *models.Project, tags []string) error {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := `
        MERGE (u:User {id: $owner_id})
        
        MERGE (p:Project {id: $id})
        SET p.title = $title,
            p.description = $description,
            p.status = $status,
            p.created_at = $created_at
        
        MERGE (u)-[r:CREATED]->(p)
        
        WITH p
        UNWIND $tags as tagName
        MERGE (t:Tag {name: tagName})
        MERGE (p)-[rt:HAS_TAG]->(t)
        
        RETURN p.id
        `
        params := map[string]interface{}{
            "id":          project.ProjectID,
            "owner_id":    project.OwnerID,
            "title":       project.Title,
            "description": project.Description,
            "status":      project.Status,
            "created_at":  project.CreatedAt.Unix(),
            "tags":        tags,
        }
        
        result, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume()
    })
    return err
}

func (r *graphRepository) GetProjectsByTag(tagName string) ([]*models.Project, error) {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := `
        MATCH (t:Tag {name: $tagName})<-[:HAS_TAG]-(p:Project)
        RETURN p.id, p.title, p.description, p.status, p.created_at
        `
        records, err := tx.Run(query, map[string]interface{}{
            "tagName": tagName,
        })
        if err != nil {
            return nil, err
        }

        var projects []*models.Project
        for records.Next() {
            record := records.Record()
            project := &models.Project{
                ProjectInfo: models.ProjectInfo{
                    ProjectID:   record.GetByIndex(0).(string),
                    Title:       record.GetByIndex(1).(string),
                    Description: record.GetByIndex(2).(string),
                    Status:      models.Status(record.GetByIndex(3).(string)),
                    CreatedAt:   time.Unix(record.GetByIndex(4).(int64), 0),
                },
            }
            projects = append(projects, project)
        }
        return projects, nil
    })

    if err != nil {
        return nil, err
    }

    return result.([]*models.Project), nil
}

func (r *graphRepository) GetUserExpertise(userID string) (map[string]string, error) {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := `
        MATCH (u:User {id: $userID})-[r:HAS_SKILL]->(t:Tag)
        RETURN t.name, r.level
        `
        records, err := tx.Run(query, map[string]interface{}{
            "userID": userID,
        })
        if err != nil {
            return nil, err
        }

        expertise := make(map[string]string)
        for records.Next() {
            record := records.Record()
            tagName := record.GetByIndex(0).(string)
            level := record.GetByIndex(1).(string)
            expertise[tagName] = level
        }
        return expertise, nil
    })

    if err != nil {
        return nil, err
    }

    return result.(map[string]string), nil
}

func (r *graphRepository) GetUsersWithTag(tagName string) ([]string, error) {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
    query := `
    MATCH (t:Tag {name: $tagName})<-[:HAS_SKILL]-(u:User)
    RETURN u.id
    `
    records, err := tx.Run(query, map[string]interface{}{
        "tagName": tagName,
    })
    if err != nil {
        return nil, err
    }

    var userIDs []string
    for records.Next() {
        record := records.Record()
        userID := record.GetByIndex(0).(string)
        userIDs = append(userIDs, userID)
    }
    return userIDs, nil
})

if err != nil {
    return nil, err
}

return result.([]string), nil

    if err != nil {
        return nil, err
    }

    return result.([]string), nil
}

func (r *graphRepository) GetProjectCollaborators(projectID string) ([]string, error) {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := `
        MATCH (p:Project {id: $projectID})<-[:CREATED]-(u:User)
        RETURN collect(u.id) as collaborators
        `
        record, err := tx.Run(query, map[string]interface{}{
            "projectID": projectID,
        })
        if err != nil {
            return nil, err
        }

        if result, err := record.Single(); err != nil {
            return []string{}, nil
        } else {
            collaborators := result.GetByIndex(0).([]interface{})
            userIDs := make([]string, len(collaborators))
            for i, c := range collaborators {
                userIDs[i] = c.(string)
            }
            return userIDs, nil
        }
    })

    if err != nil {
        return nil, err
    }

    return result.([]string), nil
}

// Rest of the existing methods remain unchanged
func (r *graphRepository) CreateTagNode(tag *models.Tag) error {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := `
        MERGE (t:Tag {name: $name})
        SET t.category = $category,
            t.usage_count = $usage_count,
            t.created_at = $created_at,
            t.updated_at = $updated_at
        `
        params := map[string]interface{}{
            "name": tag.Name,
            "category": tag.Category,
            "usage_count": tag.UsageCount,
            "created_at": tag.CreatedAt.Unix(),
            "updated_at": tag.UpdatedAt.Unix(),
        }
        result, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume()
    })
    return err
}

func (r *graphRepository) CreateTagRelationship(tag1, tag2 string, strength float64) error {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := `
        MATCH (t1:Tag {name: $tag1})
        MATCH (t2:Tag {name: $tag2})
        MERGE (t1)-[r:RELATED_TO]-(t2)
        SET r.strength = CASE
            WHEN r.strength IS NULL THEN $strength
            ELSE r.strength + $strength
            END,
            r.co_occurrences = CASE
            WHEN r.co_occurrences IS NULL THEN 1
            ELSE r.co_occurrences + 1
            END,
            r.last_updated = timestamp()
        `
        params := map[string]interface{}{
            "tag1": tag1,
            "tag2": tag2,
            "strength": strength,
        }
        result, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume()
    })
    return err
}

func (r *graphRepository) UpdateTagNode(tag *models.Tag) error {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := `
        MATCH (t:Tag {name: $name})
        SET t.category = $category,
            t.usage_count = $usage_count,
            t.updated_at = $updated_at
        `
        params := map[string]interface{}{
            "name": tag.Name,
            "category": tag.Category,
            "usage_count": tag.UsageCount,
            "updated_at": tag.UpdatedAt.Unix(),
        }
        result, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume()
    })
    return err
}

func (r *graphRepository) GetProjectTags(projectID string) ([]string, error) {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := `
        MATCH (p:Project {id: $project_id})-[:HAS_TAG]->(t:Tag)
        RETURN collect(t.name) as tags
        `
        params := map[string]interface{}{
            "project_id": projectID,
        }

        record, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }

        if result, err := record.Single(); err != nil {
            return []string{}, nil
        } else {
            tags := result.GetByIndex(0).([]interface{})
            stringTags := make([]string, len(tags))
            for i, tag := range tags {
                stringTags[i] = tag.(string)
            }
            return stringTags, nil
        }
    })

    if err != nil {
        return nil, err
    }

    return result.([]string), nil
}

func (r *graphRepository) DeleteProjectNode(projectID string) error {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := `
        MATCH (p:Project {id: $project_id})
        DETACH DELETE p
        `
        params := map[string]interface{}{
            "project_id": projectID,
        }

        result, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume()
    })
    return err
}

func (r *graphRepository) GetProjectsWithTags(projectIDs []string) (map[string][]string, error) {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := `
        MATCH (p:Project)-[:HAS_TAG]->(t:Tag)
        WHERE p.id IN $project_ids
        WITH p.id as projectID, collect(t.name) as tags
        RETURN projectID, tags
        `
        params := map[string]interface{}{
            "project_ids": projectIDs,
        }

        records, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }

        projectTags := make(map[string][]string)
        for records.Next() {
            record := records.Record()
            projectID := record.GetByIndex(0).(string)
            tags := record.GetByIndex(1).([]interface{})
            
            stringTags := make([]string, len(tags))
            for i, tag := range tags {
                stringTags[i] = tag.(string)
            }
            projectTags[projectID] = stringTags
        }
        return projectTags, nil
    })

    if err != nil {
        return nil, err
    }

    return result.(map[string][]string), nil
}

func (r *graphRepository) GetRelatedTags(tagName string, minStrength float64) ([]*models.TagRelationship, error) {
    // First get the tag ID from PostgreSQL for the source tag
    // var sourceTagID int
    // err := r.db.QueryRow(context.Background(),
    //     "SELECT id FROM tags WHERE name = $1",
    //     tagName).Scan(&sourceTagID)
    // if err != nil {
    //     return nil, fmt.Errorf("error getting source tag ID: %v", err)
    // }

    // Get related tags from Neo4j
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := `
        MATCH (t1:Tag {name: $tagName})-[r:RELATED_TO]-(t2:Tag)
        WHERE r.strength >= $minStrength
        RETURN t2.name, r.strength, r.co_occurrences, r.last_updated
        ORDER BY r.strength DESC
        `
        params := map[string]interface{}{
            "tagName": tagName,
            "minStrength": minStrength,
        }

        records, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }

        var relationships []*models.TagRelationship
        for records.Next() {
            // record := records.Record()
            // targetTagName := record.GetByIndex(0).(string)
            
            // Get target tag ID from PostgreSQL
            // var targetTagID int
            // err := r.db.QueryRow(context.Background(),
            //     "SELECT id FROM tags WHERE name = $1",
            //     targetTagName).Scan(&targetTagID)
            // if err != nil {
            //     continue // Skip if tag not found in PostgreSQL
            // }

            // relationships = append(relationships, &models.TagRelationship{
            //     Tag1ID: sourceTagID,
            //     Tag2ID: targetTagID,
            //     Strength: record.GetByIndex(1).(float64),
            //     CoOccurrences: int(record.GetByIndex(2).(int64)),
            //     LastUpdated: time.Unix(record.GetByIndex(3).(int64), 0),
            // })
        }
        return relationships, nil
    })

    if err != nil {
        return nil, err
    }

    return result.([]*models.TagRelationship), nil
}

func (r *graphRepository) SetUserTagExpertise(userID, tagName string, level string, years int) error {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := `
        MATCH (u:User {id: $userID}), (t:Tag {name: $tagName})
        MERGE (u)-[r:HAS_SKILL]->(t)
        SET r.level = $level,
            r.years = $years,
            r.updated_at = timestamp()
        `
        params := map[string]interface{}{
            "userID": userID,
            "tagName": tagName,
            "level": level,
            "years": years,
        }
        result, err := tx.Run(query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume()
    })
    return err
}

func (r *graphRepository) GetTagExperts(tagName string) ([]models.User, error) {
    session := r.driver.NewSession(neo4j.SessionConfig{})
    defer session.Close()

    result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
        query := `
        MATCH (u:User)-[r:HAS_SKILL]->(t:Tag {name: $tagName})
        WHERE r.level IN ['Expert', 'Advanced']
        RETURN u.id, u.username, r.level, r.years
        ORDER BY r.years DESC
        `
        records, err := tx.Run(query, map[string]interface{}{"tagName": tagName})
        if err != nil {
            return nil, err
        }

        var users []models.User
        for records.Next() {
            record := records.Record()
            users = append(users, models.User{
                ID: record.GetByIndex(0).(string),
                Username: record.GetByIndex(1).(string),
            })
        }
        return users, nil
    })

    if err != nil {
        return nil, err
    }

    return result.([]models.User), nil
}
