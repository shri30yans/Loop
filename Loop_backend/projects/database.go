package projects

import (
	db "Loop/database"
	. "Loop/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgtype"
)


type ErrNoProjects struct {
    Keyword string
}

func (e *ErrNoProjects) Error() string {
    return fmt.Sprintf("no projects found matching keyword: %s", e.Keyword)
}


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

func FetchProjects(keyword *string) ([]ProjectsResponse, int, error) {
	baseQuery := `
        SELECT 
            COUNT(*) OVER() AS total_projects,
            p.project_id, 
            p.owner_id, 
            p.title, 
            p.description, 
            p.status, 
            p.created_at,
            COALESCE(
                json_agg(
                    DISTINCT pt.tag_description
                ) FILTER (WHERE pt.tag_description IS NOT NULL),
                '[]'
            ) as tags
        FROM projects p
        LEFT JOIN project_tags pt ON p.project_id = pt.project_id
    `
	var query string
	var args []interface{}

	// Modify query based on keyword
	if keyword != nil {
		query = baseQuery + `
            WHERE p.title ILIKE '%' || $1 || '%'
        GROUP BY p.project_id, p.owner_id, p.title, p.description, p.status, p.created_at`
		args = append(args, *keyword)
	} else {
		query = baseQuery + `
        GROUP BY p.project_id, p.owner_id, p.title, p.description, p.status, p.created_at`
	}

	rows, err := db.DB.Query(context.Background(), query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching projects: %v", err)
	}
	defer rows.Close()

	var projects []ProjectsResponse
	var totalProjects int

	// Iterate through rows and collect project data
	for rows.Next() {
		var project ProjectsResponse
		var status pgtype.Text
		var tagsData []byte

		err := rows.Scan(
			&totalProjects,
			&project.ProjectID,
			&project.OwnerID,
			&project.Title,
			&project.Description,
			&status,
			&project.CreatedAt,
			&tagsData,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning project row: %v", err)
		}

		// Decode tags
		if err := json.Unmarshal(tagsData, &project.Tags); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling tags: %v", err)
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating project rows: %v", err)
	}

	if len(projects) == 0 {
		if keyword != nil {
			return nil, 0, &ErrNoProjects{Keyword: *keyword}
		}
		return nil, 0, &ErrNoProjects{}
	}

	return projects, totalProjects, nil
}


func FetchProjectInfo(projectID int) ([]Project, error) {
    query := `
        SELECT 
            p.project_id, 
            p.owner_id, 
            p.title, 
            p.description, 
            p.introduction, 
            p.status, 
            p.created_at,
            COALESCE(u.name, '') as owner_name,
            COALESCE(u.email, '') as owner_email,
            u.bio as owner_bio,
            u.location as owner_location,
            COALESCE(
                (SELECT json_agg(
                    DISTINCT jsonb_build_object(
                        'section_id', ps2.section_number, 
                        'title', ps2.title, 
                        'body', ps2.body
                    )
                )
                FROM project_sections ps2 
                WHERE ps2.project_id = p.project_id 
                AND ps2.section_number IS NOT NULL),
                '[]'::json
            ) AS sections,
            COALESCE(
                (SELECT json_agg(DISTINCT pt2.tag_description)
                FROM project_tags pt2 
                WHERE pt2.project_id = p.project_id 
                AND pt2.tag_description IS NOT NULL),
                '[]'::json
            ) AS tags
        FROM 
            projects p
        LEFT JOIN 
            users u ON p.owner_id = u.id
        WHERE 
            p.project_id = $1
        GROUP BY 
            p.project_id, p.owner_id, p.title, p.description, 
            p.introduction, p.status, p.created_at,
            u.name, u.email, u.bio, u.location
    `

	rows, err := db.DB.Query(context.Background(), query, projectID)
	if err != nil {
		return nil, fmt.Errorf("error fetching projects: %v", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		var sectionsData []byte
		var tagsJson string      
		var bio, location *string

		// Simple initialization
		p.Owner = UserDetails{}

		err := rows.Scan(
			&p.ProjectID,
			&p.OwnerID,
			&p.Title,
			&p.Description,
			&p.Introduction,
			&p.Status,
			&p.CreatedAt,
			&p.Owner.Name,
			&p.Owner.Email,
			&bio,
			&location,
			&sectionsData,
			&tagsJson,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning project row: %v", err)
		}

		// Assign after scanning
		p.Owner.Bio = bio
		p.Owner.Location = location
		p.Owner.ID = p.OwnerID

		// Set owner ID from project's owner_id
		p.Owner.ID = p.OwnerID
		p.Owner.CreatedAt = p.CreatedAt 

		// Decode sections
		if err := json.Unmarshal(sectionsData, &p.Sections); err != nil {
			return nil, fmt.Errorf("error unmarshaling sections: %v", err)
		}

		// Decode tags from JSON array string
		if err := json.Unmarshal([]byte(tagsJson), &p.Tags); err != nil {
			return nil, fmt.Errorf("error unmarshaling tags: %v", err)
		}

		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating project rows: %v", err)
	}

	return projects, nil
}
