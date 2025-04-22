package queries

import (
	"fmt"
	"strings"
)

// QueryManager handles all graph database queries
type QueryManager struct{}

var manager *QueryManager

// GetQueryManager returns singleton instance of QueryManager
func GetQueryManager() *QueryManager {
	if manager == nil {
		manager = &QueryManager{}
	}
	return manager
}

// FormatCreateRelationship formats the relationship creation query based on node types
func (qm *QueryManager) FormatCreateRelationship(sourceType, targetType, relationType string, useIdentityBased bool) string {
	if useIdentityBased {
		return fmt.Sprintf(IdentityBasedRelationshipQuery,
			sourceType,
			targetType,
			relationType)
	}
	return fmt.Sprintf(MixedIdentityRelationshipQuery,
		sourceType,
		targetType,
		relationType)
}

// FormatEntityQuery formats entity queries based on type
func (qm *QueryManager) FormatEntityQuery(entityType string, isIdentityBased bool) string {
	if isIdentityBased {
		return fmt.Sprintf(CREATE_IDENTITY_ENTITY_QUERY, entityType)
	}
	// For reusable entities, we need to format both the node type and relationship type
	return fmt.Sprintf(MERGE_REUSABLE_ENTITY_QUERY,
		entityType,
		strings.ToUpper(entityType),
	)
}

// GetProjectQuery returns the appropriate project query
func (qm *QueryManager) GetProjectQuery(operation string) string {
	switch operation {
	case "create":
		return CreateProjectNodeQuery
	case "link_owner":
		return LinkProjectToOwnerQuery
	case "get_nodes":
		return GetProjectNodesQuery
	case "get_relationships":
		return GetProjectRelationshipsQuery
	case "delete":
		return DeleteProjectQuery
	default:
		return ""
	}
}

// GetEntityQuery returns the appropriate entity query
func (qm *QueryManager) GetEntityQuery(operation string) string {
	switch operation {
	case "get_by_project_type":
		return GET_ENTITIES_BY_PROJECT_TYPE_QUERY
	case "get_by_type":
		return GET_ENTITIES_BY_TYPE_QUERY
	case "delete":
		return DELETE_ENTITY_QUERY
	case "batch_create":
		return BATCH_CREATE_ENTITIES_QUERY
	default:
		return ""
	}
}

// GetPersonQuery returns the appropriate person query
func (qm *QueryManager) GetPersonQuery(operation string) string {
	switch operation {
	case "create":
		return CreatePersonNodeQuery
	case "update":
		return UpdatePersonProfileQuery
	case "get_projects":
		return GetPersonProjectsQuery
	case "get_contributions":
		return GetPersonContributionsQuery
	default:
		return ""
	}
}

// FormatPersonEntityLink formats queries for linking person to entities
func (qm *QueryManager) FormatPersonEntityLink(entityType, idField, relationType string) string {
	return fmt.Sprintf(LinkPersonToEntityQuery,
		entityType,
		idField,
		relationType)
}

// FormatPersonRelatedEntities formats queries for getting person-related entities
func (qm *QueryManager) FormatPersonRelatedEntities(relationType, entityType string) string {
	return fmt.Sprintf(GetPersonRelatedEntitiesQuery,
		relationType,
		entityType)
}

// GetRelationshipQuery returns the appropriate relationship query
func (qm *QueryManager) GetRelationshipQuery(operation string) string {
	switch operation {
	case "create":
		return CreateRelationshipQuery
	case "update":
		return UpdateRelationshipQuery
	case "delete":
		return DeleteRelationshipQuery
	case "get_by_type":
		return GetRelationshipsByTypeQuery
	case "get_all":
		return GetEntityRelationshipsQuery
	case "get_properties":
		return GetRelationshipPropertiesQuery
	case "count":
		return CountRelationshipsByTypeQuery
	default:
		return ""
	}
}

// FormatBatchRelationship formats batch relationship operations
func (qm *QueryManager) FormatBatchRelationship(sourceType, targetType, relationType string) string {
	return fmt.Sprintf(BatchCreateRelationshipsQuery,
		sourceType,
		targetType,
		relationType)
}

// FormatCheckRelationship formats the relationship check query
func (qm *QueryManager) FormatCheckRelationship(sourceType, targetType, relationType string, useID bool) string {
	idField := "id"
	if !useID {
		idField = "name"
	}
	return fmt.Sprintf(CheckRelationshipExistsQuery,
		sourceType,
		idField,
		targetType,
		idField,
		relationType)
}

// FormatPathQuery formats the path finding query with max depth
func (qm *QueryManager) FormatPathQuery(maxDepth int) string {
	return fmt.Sprintf(FindPathsBetweenEntitiesQuery, maxDepth)
}
