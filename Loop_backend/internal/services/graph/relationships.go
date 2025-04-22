package graph

import (
    "strings"
    "Loop_backend/platform/database/neo4j/entities"
)
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
    entities.TypeTechnology:  Development,
    entities.TypeFeature:     Implements,
    entities.TypeTag:         RelatedTo,
    entities.TypeCategory:    BelongsTo,
    entities.TypeStakeholder: HasStakeholder,
    entities.TypePerson:      DevelopedBy,
    entities.TypeMethodology: Uses,
}

// GetRelationshipType returns the appropriate relationship type between two entities
func GetRelationshipType(sourceType, targetType string) string {
    if !strings.EqualFold(sourceType, entities.TypeProject) {
        return RelatedTo
    }
    
    for key, relation := range ProjectRelationshipMap {
        if strings.EqualFold(key, targetType) {
            return relation
        }
    }
    return RelatedTo
}
