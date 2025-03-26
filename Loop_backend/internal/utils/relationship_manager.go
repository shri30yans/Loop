package utils

import (
	"Loop_backend/platform/database/neo4j/entities"
	"strings"
	"sync"
)

type RelationshipManager struct{}

var (
	relationshipManager *RelationshipManager
	managerOnce         sync.Once
)

// GetRelationshipManager returns a singleton instance of RelationshipManager
func GetRelationshipManager() (*RelationshipManager, error) {
	managerOnce.Do(func() {
		relationshipManager = &RelationshipManager{}
	})
	return relationshipManager, nil
}

// GetRelationType determines the appropriate relationship type between two entities
func (rm *RelationshipManager) GetRelationType(sourceType, targetType string) string {
	sourceType = strings.ToLower(sourceType)
	targetType = strings.ToLower(targetType)

	return entities.GetRelationshipType(sourceType, targetType)
}

// ValidateEntityType checks if an entity type is valid
func (rm *RelationshipManager) ValidateEntityType(entityType string) bool {
	return entities.IsValidEntityType(entityType)
}
