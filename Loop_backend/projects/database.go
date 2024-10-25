package projects

import (
	db "Loop/database"
	"context"
	"fmt"
	"github.com/google/uuid"
)

func CreateProject(title, description, introduction, tags string) (uuid.UUID, error) {
	var projectID uuid.UUID
	err := db.DB.QueryRow(context.Background(),
		"INSERT INTO projects (title, description, introduction, tags) VALUES ($1, $2, $3, $4) RETURNING project_id",
		title, description, introduction, tags).Scan(&projectID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error creating project: %v", err)
	}
	fmt.Println("Created project", projectID)
	return projectID, nil
}

func FetchProjects() ([]Project, error) {
	rows, err := db.DB.Query(context.Background(), "SELECT project_id, owner_id, title, description, status, created_at FROM projects")
	if err != nil {
		return nil, fmt.Errorf("error fetching projects: %v", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		err := rows.Scan(&p.ProjectID, &p.OwnerID, &p.Title, &p.Description, &p.Status, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning project row: %v", err)
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating project rows: %v", err)
	}

	return projects, nil
}
