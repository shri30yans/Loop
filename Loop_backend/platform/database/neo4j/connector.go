package neo4j

import (
	"Loop_backend/config"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

var graphDB neo4j.Driver

// InitNeo4j initializes the graph database connection
func InitNeo4j(cfg *config.Neo4jConfig) error {
	fmt.Println("Initializing graph database...")
	var err error

	graphDB, err = neo4j.NewDriver(cfg.URI, neo4j.BasicAuth(cfg.Username, cfg.Password, ""))
	if err != nil {
		return fmt.Errorf("unable to connect to graph database: %v", err)
	}

	// Test the connection
	err = graphDB.VerifyConnectivity()
	if err != nil {
		return fmt.Errorf("unable to verify graph database connection: %v", err)
	}

	fmt.Println("Connected to Graph Database")

	// Run migrations after successful connection
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	return nil
}

// GetDriver returns the graph database driver
func GetDriver() neo4j.Driver {
	return graphDB
}

// runMigrations reads and executes all .cyp migration files
func runMigrations() error {
	migrationsDir := "platform/database/neo4j/migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	session := graphDB.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".cyp") {
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationsDir, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", file.Name(), err)
		}

		// Split the content into individual statements while preserving comments
		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" || strings.HasPrefix(stmt, "//") {
				continue
			}

			// Execute each non-comment statement
			_, err = session.Run(stmt, nil)
			if err != nil {
				return fmt.Errorf("failed to execute statement from %s: %v", file.Name(), err)
			}
		}

		fmt.Printf("Executed migration: %s\n", file.Name())
	}

	return nil
}

// Close closes the graph database driver
func Close() {
	if graphDB != nil {
		graphDB.Close()
	}
}
