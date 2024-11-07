package database

import (
	"context"
	"fmt"
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
