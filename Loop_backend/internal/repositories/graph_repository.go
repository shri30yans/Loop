package repositories

import (
	"Loop_backend/internal/models"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type GraphRepository interface {
	CreateProjectNode(project *models.Project) error
	CreateUserNode(user *models.User) error
	CreateProjectUserRelation(projectID, userID string) error
	CreateProjectTagRelations(projectID string, tags []string) error
	GetProjectTags(projectID string) ([]string, error)
	GetProjectsWithTags(projectIDs []string) (map[string][]string, error)
	DeleteProjectNode(projectID string) error
}

type graphRepository struct {
	driver neo4j.Driver
}

func NewGraphRepository(uri, username, password string) (GraphRepository, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("error creating neo4j driver: %v", err)
	}
	return &graphRepository{driver: driver}, nil
}

func (r *graphRepository) CreateProjectNode(project *models.Project) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Create project node and create the user-project relationship in a single query
		query := `
            MERGE (u:User {id: $owner_id})
            MERGE (p:Project {id: $id})
            SET p.title = $title,
                p.description = $description,
                p.status = $status,
                p.created_at = $created_at
            MERGE (u)-[r:CREATED]->(p)
            RETURN p.id
        `
		params := map[string]interface{}{
			"id":          project.ProjectID,
			"owner_id":    project.OwnerID,
			"title":       project.Title,
			"description": project.Description,
			"status":      project.Status,
			"created_at":  project.CreatedAt.Unix(),
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

		result, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}

		record, err := result.Single()
		if err != nil {
			return []string{}, nil // Return empty array if no tags found
		}

		tags := record.Values[0].([]interface{})
		stringTags := make([]string, len(tags))
		for i, tag := range tags {
			stringTags[i] = tag.(string)
		}
		return stringTags, nil
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
            RETURN p.id as project_id, collect(t.name) as tags
        `
		params := map[string]interface{}{
			"project_ids": projectIDs,
		}

		result, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}

		projectTags := make(map[string][]string)
		for result.Next() {
			record := result.Record()
			projectID := record.Values[0].(string)
			tags := record.Values[1].([]interface{})

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

func (r *graphRepository) CreateUserNode(user *models.User) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `
            MERGE (u:User {id: $id})
            SET u.username = $username,
                u.email = $email,
                u.created_at = $created_at
        `
		params := map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt.Unix(),
		}

		result, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}
		return result.Consume()
	})

	return err
}

func (r *graphRepository) CreateProjectUserRelation(projectID, userID string) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `
            MATCH (p:Project {id: $project_id})
            MATCH (u:User {id: $user_id})
            MERGE (u)-[r:CREATED]->(p)
        `
		params := map[string]interface{}{
			"project_id": projectID,
			"user_id":    userID,
		}

		result, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}
		return result.Consume()
	})

	return err
}

func (r *graphRepository) CreateProjectTagRelations(projectID string, tags []string) error {
	session := r.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `
            MATCH (p:Project {id: $project_id})
            UNWIND $tags as tag
            MERGE (t:Tag {name: tag})
            MERGE (p)-[r:HAS_TAG]->(t)
        `
		params := map[string]interface{}{
			"project_id": projectID,
			"tags":       tags,
		}

		result, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}
		return result.Consume()
	})

	return err
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
