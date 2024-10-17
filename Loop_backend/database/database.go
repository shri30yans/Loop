package database

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func InitDB() error {
	fmt.Println("Initializing database...")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	DB, err = pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	fmt.Println("Connected to Database")

	err = createTables()
	if err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	return nil
}

func createTables() error {
	fmt.Println("Creating tables")
	queries := []string{
		DropAllTables,
		CreateUsersTable,
		CreateProjectsTable,
		CreateFeedbackTable,
		CreateProjectUpdatesTable,
		CreateEventsTable,
		CreateUserEventParticipationTable,
	}

	for _, query := range queries {
		_, err := DB.Exec(context.Background(), query)
		if err != nil {
			return err
		}
	}

	fmt.Println("Created tables")

	return nil
}

func CreateProject(title, description, introduction, tags string) (uuid.UUID, error) {
	var projectID uuid.UUID
	err := DB.QueryRow(context.Background(),
		"INSERT INTO projects (title, description, introduction, tags) VALUES ($1, $2, $3, $4) RETURNING project_id",
		title, description, introduction, tags).Scan(&projectID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error creating project: %v", err)
	}
	fmt.Println("Created project", projectID)
	return projectID, nil
}

func FetchProjects() ([]Project, error) {
	rows, err := DB.Query(context.Background(), "SELECT project_id, owner_id, title, description, status, created_at FROM projects")
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
