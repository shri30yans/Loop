package entities

import "strings"

// Relationship constants
const (
	RelatedTo      = "RELATED_TO"
	Uses           = "USES"
	Implements     = "IMPLEMENTS"
	BelongsTo      = "BELONGS_TO"
	HasStakeholder = "HAS_STAKEHOLDER"
	DevelopedBy    = "DEVELOPED_BY"
	Development    = "DEVELOPMENT"
	Support        = "SUPPORT"
	Trade          = "TRADE"
)

var ProjectRelationshipMap = map[string]string{
	TypeTechnology:  Development,
	TypeFeature:     Implements,
	TypeTag:         RelatedTo,
	TypeCategory:    BelongsTo,
	TypeStakeholder: HasStakeholder,
	TypePerson:      DevelopedBy,
	TypeMethodology: Uses,
}

// Updated GetRelationshipType in relationships.go
// GetRelationshipType determines the appropriate relationship type between two entities
func GetRelationshipType(sourceType, targetType string) string {
    // Normalize types to ensure consistent matching
    sourceType = strings.TrimSpace(sourceType)
    targetType = strings.TrimSpace(targetType)
    
    // Standardize types to match our canonical types
    sourceType = normalizeType(sourceType)
    targetType = normalizeType(targetType)

    // Handle Project relationships first - Project should generally be the source
    if sourceType == TypeProject {
        if relation, ok := ProjectRelationshipMap[targetType]; ok {
            return relation
        }
    }

    // Specific rules for non-Project relationships
    switch {
    case sourceType == TypeTechnology && targetType == TypeProject:
        return Development
    case sourceType == TypeStakeholder && targetType == TypePerson:
        return DevelopedBy
    case sourceType == TypeFeature && targetType == TypeTechnology:
        return Implements
    case sourceType == TypeProject && targetType == TypePlatform:
        return Uses
    case sourceType == TypeTechnology && targetType == TypeFeature:
        return Support
    }

    // Default fallback
    return RelatedTo
}

// normalizeType ensures consistent type matching
func normalizeType(entityType string) string {
    entityType = strings.TrimSpace(entityType)
    
    // Check for direct match with existing types (case-insensitive)
    for _, validType := range EntityTypes() {
        if strings.EqualFold(entityType, validType) {
            return validType
        }
    }
    
    // Handle common synonyms
    switch strings.ToLower(entityType) {
    case "tech", "technology", "technologies":
        return TypeTechnology
    case "person", "individual", "people", "user":
        return TypePerson
    case "stakeholder", "stakeholders", "client", "customer":
        return TypeStakeholder
    case "project", "application", "app", "system", "product":
        return TypeProject
    case "feature", "functionality", "function", "capability":
        return TypeFeature
    case "tag", "label", "keyword":
        return TypeTag
    case "category", "group", "classification":
        return TypeCategory
    case "methodology", "method", "framework", "approach":
        return TypeMethodology
    case "platform", "infrastructure", "server", "service":
        return TypePlatform
    }
    
    return TypeTag // Default fallback
}
