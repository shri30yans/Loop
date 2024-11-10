package projects

import (
	db "Loop/database"
	"context"
	"fmt"
)

func CreateProject(title, description, introduction, tags string, ownerID int, sections []ProjectSection) (int, error) {
	var projectID int
	err := db.DB.QueryRow(context.Background(),
		"INSERT INTO projects (title, owner_id, description, introduction, tags) VALUES ($1, $2, $3, $4, $5) RETURNING project_id",
		title, ownerID, description, introduction, tags).Scan(&projectID)
	if err != nil {
		return 0, fmt.Errorf("error creating project: %v", err)
	}
	CreateProjectSections(projectID, sections)
	fmt.Println("Created project", projectID)
	return projectID, nil
}

func CreateProjectSections(projectID int, sections []ProjectSection) error {
	for _, section := range sections {
		fmt.Println("Creating project section", section.UpdateNumber)
		_, err := db.DB.Exec(context.Background(),
			"INSERT INTO project_sections (section_id, project_id, title, body) VALUES ($1, $2, $3, $4)", section.UpdateNumber, projectID, section.Title,section.Body)
		if err != nil {
			return fmt.Errorf("error creating project section: %v", err)
		}
	}
	return nil
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
