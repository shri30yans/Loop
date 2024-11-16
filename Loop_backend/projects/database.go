package projects

import (
	db "Loop/database"
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgtype"
)

func CreateProject(title, description, introduction string, tags []string, ownerID int, sections []ProjectSection) (int, error) {
	var projectID int

	sectionsJSON, err := json.Marshal(sections)
	if err != nil {
		return 0, fmt.Errorf("error marshaling sections: %v", err)
	}

	query := `SELECT create_project($1, $2, $3, $4, $5::text[], $6::jsonb)`

	tagsArray := pgtype.TextArray{}
	if err := tagsArray.Set(tags); err != nil {
		return 0, fmt.Errorf("error converting tags to array: %v", err)
	}

	err = db.DB.QueryRow(
		context.Background(),
		query,
		title, description, introduction, ownerID, tagsArray, sectionsJSON,
	).Scan(&projectID)
	if err != nil {
		return 0, fmt.Errorf("error creating project: %v", err)
	}

	fmt.Println("Created project with ID:", projectID)
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
		var project Project
		var status pgtype.Text
		err := rows.Scan(&project.ProjectID, &project.OwnerID, &project.Title, &project.Description, &status, &project.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning project row: %v", err)
		}
		projects = append(projects, project)
	}

	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating project rows: %v", err)
	}

	return projects, nil
}

func FetchProjectInfo(projectID int) ([]Project, error) {
	fmt.Println("Project fetched for", projectID)

	query := `
		SELECT 
			p.project_id, 
			p.owner_id, 
			p.title, 
			p.description, 
			p.introduction, 
			p.status, 
			p.created_at, 
			COALESCE(
				json_agg(
					DISTINCT jsonb_build_object(
						'section_id', ps.section_number, 
						'title', ps.title, 
						'body', ps.body
					)
				) FILTER (WHERE ps.section_number IS NOT NULL), 
				'[]'
			) AS sections,
			COALESCE(
				array_agg(
					DISTINCT pt.tag_description
				) FILTER (WHERE pt.tag_description IS NOT NULL),
				'{}'
			) AS tags
		FROM 
			projects p
		LEFT JOIN 
			project_sections ps 
		ON 
			p.project_id = ps.project_id
		LEFT JOIN 
			project_tags pt 
		ON 
			p.project_id = pt.project_id
		WHERE 
			p.project_id = $1
		GROUP BY 
			p.project_id, p.owner_id, p.title, p.description, p.introduction, p.status, p.created_at
	`
	// Query explaination:
	// json_agg: Aggregates all sections into a JSON array grouped by project_id.
	// jsonb_build_object: Constructs JSON objects for individual sections.
	// FILTER: Ensures NULL values for sections do not disrupt the JSON array aggregation.
	// COALESCE: Defaults the sections array to an empty array ('[]') if no sections exist.

	rows, err := db.DB.Query(context.Background(), query, projectID)
	if err != nil {
		return nil, fmt.Errorf("error fetching projects: %v", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		var sectionsData []byte
		var tagsData []string

		// Scan project data, sections JSON, and tags array
		err := rows.Scan(&p.ProjectID, &p.OwnerID, &p.Title, &p.Description, &p.Introduction, &p.Status, &p.CreatedAt, &sectionsData, &tagsData)
		if err != nil {
			return nil, fmt.Errorf("error scanning project row: %v", err)
		}

		// Decode the JSON array of sections into the Project.Sections field
		if err := json.Unmarshal(sectionsData, &p.Sections); err != nil {
			return nil, fmt.Errorf("error unmarshaling sections: %v", err)
		}

		// Assign tags directly to the Project.Tags field
		p.Tags = tagsData

		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating project rows: %v", err)
	}

	return projects, nil
}
