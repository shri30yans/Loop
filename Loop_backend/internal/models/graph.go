package models

// KnowledgeGraph represents a structured graph of entities and their relationships
type KnowledgeGraph struct {
	Entities      []Entity
	Relationships []Relationship
	Keywords      []string
}

// Entity represents a node in the knowledge graph
type Entity struct {
	Name        string
	Type        string
	Description string
	Properties  map[string]interface{}
}

// Relationship represents a connection between two entities
type Relationship struct {
	Source      string
	Target      string
	Type        string
	Description string
	Weight      int
	Category    string
}
