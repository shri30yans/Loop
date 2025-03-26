package repositories

import (
	"Loop_backend/internal/models"
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"time"
)

type GraphRepository interface {
	// Knowledge Graph operations
	StoreProjectGraph(projectID uuid.UUID, ownerID string, graph *models.KnowledgeGraph) error
	GetProjectGraph(projectID uuid.UUID) (*models.KnowledgeGraph, error)

	GetProjectTags(projectID string) ([]string, error)
	GetProjectsWithTags(projectIDs []string) (map[string][]string, error)
	DeleteProjectNode(projectID string) error

	// Tag operations
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

// StoreProjectGraph stores the entire knowledge graph for a project
func (r *graphRepository) StoreProjectGraph(projectID uuid.UUID, ownerID string, graph *models.KnowledgeGraph) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// First create project owner as Person node and create DEVELOPED_BY relationship with Project
		_, err := tx.Run(`
            MERGE (person:Person {id: $owner_id})
            SET person.role = 'owner'
            WITH person
            MERGE (project:Project {id: $project_id, project_id: $project_id})
            MERGE (project)-[r:DEVELOPED_BY]->(person)
            SET r.description = 'Project developer',
                r.weight = 10,
                r.category = 'ownership'
            `,
			map[string]interface{}{
				"owner_id":   ownerID,
				"project_id": projectID.String(),
			})
		if err != nil {
			return nil, err
		}

		// Then store all other entities with their specific node labels
		for _, entity := range graph.Entities {
			// Create the node with its specific label based on type
			query := fmt.Sprintf(`
                MERGE (n:%s {
                    project_id: $project_id,
                    name: $name
                })
                SET n.description = $description
                `, entity.Type)

			_, err := tx.Run(query,
				map[string]interface{}{
					"project_id":  projectID.String(),
					"name":        entity.Name,
					"description": entity.Description,
				})
			if err != nil {
				return nil, err
			}
		}

		// Store relationships between specifically labeled nodes
		for _, rel := range graph.Relationships {
			// Find source and target entity types
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

			query := fmt.Sprintf(`
                MATCH (n1:%s {project_id: $project_id, name: $source})
                MATCH (n2:%s {project_id: $project_id, name: $target})
                MERGE (n1)-[r:%s]->(n2)
                SET r.description = $description,
                    r.weight = $weight,
                    r.category = $category
                `, sourceType, targetType, rel.Type)

			_, err := tx.Run(query,
				map[string]interface{}{
					"project_id":  projectID.String(),
					"source":      rel.Source,
					"target":      rel.Target,
					"description": rel.Description,
					"weight":      rel.Weight,
					"category":    rel.Category,
				})
			if err != nil {
				return nil, err
			}
		}

		// Store keywords as tags
		for _, keyword := range graph.Keywords {
			_, err := tx.Run(`
                MERGE (t:Tag {name: $keyword})
                WITH t
                MATCH (p:Project {id: $project_id})
                MERGE (p)-[r:HAS_TAG]->(t)
                `,
				map[string]interface{}{
					"project_id": projectID.String(),
					"keyword":    keyword,
				})
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	return err
}

// GetProjectGraph retrieves the entire knowledge graph for a project
func (r *graphRepository) GetProjectGraph(projectID uuid.UUID) (*models.KnowledgeGraph, error) {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Get all nodes with project_id property
		entitiesResult, err := tx.Run(`
            MATCH (n)
            WHERE n.project_id = $project_id
            RETURN n.name, labels(n)[0], n.description
            `,
			map[string]interface{}{
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

		// Get all relationships between nodes with project_id
		relationshipsResult, err := tx.Run(`
            MATCH (n1)-[r]->(n2)
            WHERE n1.project_id = $project_id AND n2.project_id = $project_id
            RETURN n1.name, n2.name, type(r), r.description, r.weight, r.category
            `,
			map[string]interface{}{
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

		// Get keywords (stored as tags)
		keywordsResult, err := tx.Run(`
            MATCH (p:Project {id: $project_id})-[:HAS_TAG]->(t:Tag)
            RETURN collect(t.name) as keywords
            `,
			map[string]interface{}{
				"project_id": projectID.String(),
			})
		if err != nil {
			return nil, err
		}

		var keywords []string
		if keywordsResult.Next() {
			record := keywordsResult.Record()
			keywordsList := record.GetByIndex(0).([]interface{})
			keywords = make([]string, len(keywordsList))
			for i, k := range keywordsList {
				keywords[i] = k.(string)
			}
		}

		return &models.KnowledgeGraph{
			Entities:      entities,
			Relationships: relationships,
			Keywords:      keywords,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.KnowledgeGraph), nil
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
			projectID, _ := uuid.Parse(record.Values[0].(string))
			project := &models.Project{
				ProjectInfo: models.ProjectInfo{
					ProjectID:   projectID,
					Title:       record.Values[1].(string),
					Description: record.Values[2].(string),
					Status:      models.Status(record.Values[3].(string)),
					CreatedAt:   time.Unix(record.Values[4].(int64), 0),
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
			"name":        tag.Name,
			"category":    tag.Category,
			"usage_count": tag.UsageCount,
			"created_at":  tag.CreatedAt.Unix(),
			"updated_at":  tag.UpdatedAt.Unix(),
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
			"name":        tag.Name,
			"category":    tag.Category,
			"usage_count": tag.UsageCount,
			"updated_at":  tag.UpdatedAt.Unix(),
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
			"tag1":     tag1,
			"tag2":     tag2,
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
				ID:       record.GetByIndex(0).(string),
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
			"userID":  userID,
			"tagName": tagName,
			"level":   level,
			"years":   years,
		}
		result, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}
		return result.Consume()
	})
	return err
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
        RETURN collect(u.id) as users
        `
		record, err := tx.Run(query, map[string]interface{}{
			"tagName": tagName,
		})
		if err != nil {
			return nil, err
		}

		if result, err := record.Single(); err != nil {
			return []string{}, nil
		} else {
			users := result.GetByIndex(0).([]interface{})
			userIDs := make([]string, len(users))
			for i, u := range users {
				userIDs[i] = u.(string)
			}
			return userIDs, nil
		}
	})

	if err != nil {
		return nil, err
	}

	return result.([]string), nil
}

func (r *graphRepository) GetRelatedTags(tagName string, minStrength float64) ([]*models.TagRelationship, error) {
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
			"tagName":     tagName,
			"minStrength": minStrength,
		}

		records, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}

		var relationships []*models.TagRelationship
		for records.Next() {
			record := records.Record()
			targetTagName := record.GetByIndex(0).(string)
			relationships = append(relationships, &models.TagRelationship{
				Tag2Name:      targetTagName,
				Strength:      record.GetByIndex(1).(float64),
				CoOccurrences: int(record.GetByIndex(2).(int64)),
				UpdatedAt:     time.Unix(record.GetByIndex(3).(int64), 0),
			})
		}
		return relationships, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]*models.TagRelationship), nil
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
