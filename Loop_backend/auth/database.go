package auth

import (
	db "Loop/database"
	. "Loop/models"
	"context"
	"errors"
	"strings"
)

func CreateUser(name string, email string, hashedPassword string) (User, error) {
	var user User
	err := db.DB.QueryRow(
		context.Background(),
		"INSERT INTO users (name, email, hashed_password) VALUES ($1, $2, $3) RETURNING id,name,email, hashed_password",
		name, email, hashedPassword,
	).Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return User{}, errors.New("email already exists")
		}
		return User{}, err
	}
	return user, nil
}

func GetUserInfoById(id string) (UserInfoSummary, error) {
	var user UserInfoSummary
	user.Projects = make([]ProjectsResponse, 0)

	err := db.DB.QueryRow(
		context.Background(),
		`SELECT 
			u.id,
			u.name, 
			u.email, 
			u.location, 
			u.bio, 
			u.created_at,
			(
				SELECT json_agg(json_build_object(
					'project_id', p.project_id,
					'owner_id', p.owner_id,
					'title', p.title,
					'description', p.description,
					'introduction', p.introduction,
					'status', p.status,
					'created_at', to_char(p.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
					'sections', (
						SELECT COALESCE(
							json_agg(
								json_build_object(
									'section_id', section_number,
									'title', title,
									'body', body
								)
							),
							'[]'
						)
						FROM project_sections
						WHERE project_id = p.project_id
					),
					'tags', (
						SELECT COALESCE(
							array_agg(DISTINCT tag_description),
							'{}'
						)
						FROM project_tags
						WHERE project_id = p.project_id
					)
				))
				FROM projects p 
				WHERE p.owner_id = u.id
			) as projects
		FROM users u 
		WHERE u.id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Location,
		&user.Bio,
		&user.CreatedAt,
		&user.Projects,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetUserByEmail(email string) (User, error) {
	var user User
	err := db.DB.QueryRow(
		context.Background(),
		"SELECT id, email, hashed_password, created_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.HashedPassword, &user.CreatedAt)
	return user, err
}

