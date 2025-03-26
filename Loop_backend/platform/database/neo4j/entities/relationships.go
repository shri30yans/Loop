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

// GetRelationshipType returns the appropriate relationship type between two entities
func GetRelationshipType(sourceType, targetType string) string {
    if !strings.EqualFold(sourceType, TypeProject) {
        return RelatedTo
    }
    
    for key, relation := range ProjectRelationshipMap {
        if strings.EqualFold(key, targetType) {
            return relation
        }
    }
    return RelatedTo
}
