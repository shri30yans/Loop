package neo4j

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"Loop_backend/config"
)

var graphDB neo4j.Driver

// InitGraph initializes the graph database connection
func InitGraph(cfg *config.Neo4jConfig) error {
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
	return nil
}

// GetDriver returns the graph database driver
func GetDriver() neo4j.Driver {
	return graphDB
}

// Close closes the graph database driver
func Close() {
	if graphDB != nil {
		graphDB.Close()
	}
}
