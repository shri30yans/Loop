package models

// KnowledgeGraph represents a project's knowledge graph
type KnowledgeGraph struct {
    Name          string         `json:"name"`
    Description   string         `json:"description"`
    Relationships []Relationship `json:"relationships"`
    Entities      []Entity      `json:"entities"`
    Keywords      []string      `json:"keywords"`
}

// Relationship represents a connection between entities
type Relationship struct {
    Source string `json:"source"`
    Target string `json:"target"`
    Type   string `json:"type"`
    Weight int    `json:"weight"`
}

// Entity represents a node in the knowledge graph
type Entity struct {
    Name         string `json:"name"`
    Type         string `json:"type"`
    Level        string `json:"level"`
    Description  string `json:"description"`
    ParentEntity string `json:"parent_entity,omitempty"`
}

// Common entity levels
const (
    LevelMain = "MAIN"
    LevelSub  = "SUB"
)

// Entity types
const (
    TypeTechnology  = "TECHNOLOGY"
    TypeDomain      = "DOMAIN"
    TypeFeature     = "FEATURE"
    TypePerson      = "PERSON"
)

// Relationship types
const (
    RelationImplements   = "IMPLEMENTS"
    RelationUses        = "USES"
    RelationDependsOn   = "DEPENDS_ON"
    RelationContributes = "CONTRIBUTES"
    RelationSimilarTo   = "SIMILAR_TO"
)

// Weight ranges
const (
    WeightWeak   = 1
    WeightMedium = 5
    WeightStrong = 10
)
