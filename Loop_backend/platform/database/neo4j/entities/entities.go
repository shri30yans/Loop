package entities

import (
	"strings"
)

// Entity Types
const (
	TypeProject     = "Project"
	TypeTechnology  = "Technology"
	TypeFeature     = "Feature"
	TypeTag         = "Tag"
	TypeCategory    = "Category"
	TypeStakeholder = "Stakeholder"
	TypePerson      = "Person"
	TypePlatform    = "Platform"
	TypeMethodology = "Methodology"
)

// EntityTypes returns a slice of all valid entity types
func EntityTypes() []string {
	return []string{
		TypeProject,
		TypeTechnology,
		TypeFeature,
		TypeTag,
		TypeCategory,
		TypeStakeholder,
		TypePerson,
		TypePlatform,
		TypeMethodology,
	}
}

// IsValidEntityType checks if the given type is a valid entity type
func IsValidEntityType(entityType string) bool {
	for _, validType := range EntityTypes() {
		if strings.EqualFold(validType, entityType) {
			return true
		}
	}
	return false
}
