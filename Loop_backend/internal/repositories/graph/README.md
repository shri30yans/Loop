# Graph Repository

This package implements the graph database operations using Neo4j.

## Structure

```
graph/
├── impl/                           # Implementation files
│   ├── base.go                    # Common interfaces and base repo
│   ├── graph_entity_repo.go       # Entity operations 
│   ├── graph_project_repo.go      # Project operations
│   ├── graph_relationship_repo.go  # Relationship operations
│   └── graph_person_repo.go       # Person-specific operations
├── queries/                       # Query definitions
│   ├── entity_queries.go         # Entity-related queries
│   ├── project_queries.go        # Project-related queries
│   ├── relationship_queries.go   # Relationship-related queries
│   ├── user_queries.go          # Person/user-related queries
│   └── query_manager.go         # Query formatting and management
└── graph_repository.go           # Main repository implementation

## Usage

The graph repository is designed using the composite pattern, where each specialized repository handles its specific domain operations:

- EntityRepository: Handles generic entity CRUD operations
- ProjectRepository: Manages project graphs and their relationships
- RelationshipRepository: Manages relationships between entities
- PersonRepository: Handles person/user specific operations

### Example Usage

```go
// Create a new graph repository
repo := graph.NewGraphRepository(driver)

// Store a project graph
err := repo.StoreProjectGraph(projectID, ownerID, graph)

// Get entities by type
entities, err := repo.GetEntitiesByType("Technology")

// Create relationships
err := repo.CreateRelationship("sourceID", "targetID", "USES", props)
```

## Implementation Details

The repository is split into focused implementations to improve maintainability:

1. Each domain has its own repository interface and implementation
2. Query definitions are separated from their usage
3. The main repository composes all specialized repositories
4. Common functionality is shared through the base repository

## Deprecated Files

- graph_user_repo.go: Replaced by impl/graph_person_repo.go
- queries/tag_queries.go: Tags are now handled as standard entities
