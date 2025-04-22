package queries

// Relationship Operations
const (
	CreateRelationshipQuery = `
        MATCH (source:%s {%s: $source_id})
        MATCH (target:%s {%s: $target_id})
        MERGE (source)-[r:%s]->(target)
        SET r += $props
    `

	UpdateRelationshipQuery = `
        MATCH (source)-[r:%s]->(target)
        WHERE ID(r) = $relationship_id
        SET r += $props
    `

	DeleteRelationshipQuery = `
        MATCH ()-[r]->()
        WHERE ID(r) = $relationship_id
        DELETE r
    `

	GetRelationshipsByTypeQuery = `
        MATCH (source)-[r:%s]->(target)
        WHERE source.project_id = $source_id
        RETURN target.project_id, type(r), properties(r)
    `

	GetEntityRelationshipsQuery = `
        MATCH (n {project_id: $entity_id})-[r]-(m)
        RETURN type(r), 
               CASE WHEN startNode(r) = n THEN 'outgoing' ELSE 'incoming' END as direction,
               CASE WHEN startNode(r) = n THEN m.project_id ELSE startNode(r).project_id END as connected_id,
               properties(r) as properties
    `

	// Standard relationship queries by entity type
	IdentityBasedRelationshipQuery = `
        MATCH (n1:%s {project_id: $source_id})
        MATCH (n2:%s {project_id: $target_id})
        MERGE (n1)-[r:%s]->(n2)
        SET r.description = $description,
            r.weight = $weight,
            r.category = $category
    `

	MixedIdentityRelationshipQuery = `
        MATCH (n1:%s)
        WHERE (n1.project_id = $source_id OR n1.name = $source)
        MATCH (n2:%s)
        WHERE (n2.project_id = $target_id OR n2.name = $target)
        MERGE (n1)-[r:%s]->(n2)
        SET r.description = $description,
            r.weight = $weight,
            r.category = $category
    `

	// Batch operations
	BatchCreateRelationshipsQuery = `
        UNWIND $relationships as rel
        MATCH (source:%s {project_id: rel.source_id})
        MATCH (target:%s {project_id: rel.target_id})
        MERGE (source)-[r:%s]->(target)
        SET r += rel.properties
    `

	// Utility queries
	GetRelationshipPropertiesQuery = `
        MATCH (source)-[r:%s]->(target)
        WHERE source.project_id = $source_id AND target.project_id = $target_id
        RETURN properties(r)
    `

	CountRelationshipsByTypeQuery = `
        MATCH (n {project_id: $entity_id})-[r:%s]-()
        RETURN count(r) as count
    `

	CheckRelationshipExistsQuery = `
        MATCH (source:%s {%s: $source_id})
        MATCH (target:%s {%s: $target_id})
        RETURN EXISTS((source)-[:%s]->(target)) as exists
    `

	FindPathsBetweenEntitiesQuery = `
        MATCH path = shortestPath((source {project_id: $source_id})-[*..%d]->(target {project_id: $target_id}))
        RETURN path, length(path) as path_length
        ORDER BY path_length
        LIMIT $limit
    `
)
