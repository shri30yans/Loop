package postgres

import (
	"Loop_backend/config"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool

// InitDB initializes the database connection and runs migrations
func InitDB(cfg *config.RelationalDatabaseConfig) error {
	fmt.Println("Initializing database...")
	connString := cfg.ConnectionString
	var err error
	db, err = pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return fmt.Errorf("unable to ping database: %v", err)
	}

	fmt.Println("Connected to Database")

	err = runMigrations()
	if err != nil {
		return fmt.Errorf("error running migrations: %v", err)
	}

	return nil
}

// GetDB returns the database connection pool
func GetDB() *pgxpool.Pool {
	return db
}

// Close closes the database connection pool
func Close() {
	if db != nil {
		db.Close()
	}
}

// runMigrations executes all SQL migration files in order
func runMigrations() error {
	migrationPath := "./platform/database/postgres/migrations"
	fmt.Println("Running migrations...")
	files, err := os.ReadDir(migrationPath)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %v", err)
	}

	// Get all SQL files and sort them
	var sqlFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Execute migrations in order
	for _, fileName := range sqlFiles {
		//fmt.Printf("Executing migration: %s\n", fileName)
		query, err := os.ReadFile(filepath.Join(migrationPath, fileName))
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %v", fileName, err)
		}

		_, err = db.Exec(context.Background(), string(query))
		if err != nil {
			return fmt.Errorf("error executing migration file %s: %v", fileName, err)
		}
		fmt.Printf("Successfully executed migration: %s\n", fileName)
	}

	fmt.Println("All migrations completed successfully")
	return nil
}
