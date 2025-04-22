package queries

// Entity Operations
const (
	// BATCH_CREATE_ENTITIES_QUERY is used for bulk entity creation
	BATCH_CREATE_ENTITIES_QUERY = `
        UNWIND $entities as entity
        WITH entity, 
             CASE 
                WHEN entity.identity = true THEN {name: entity.name}
                ELSE {project_id: $project_id, name: entity.name}
             END as properties
        MERGE (n:` + "`%s`" + ` {name: entity.name})
        ON CREATE SET
            n.id = coalesce(entity.id, randomUUID()),
            n.description = entity.description,
            n.created_at = datetime()
        ON MATCH SET
            n.description = coalesce(entity.description, n.description),
            n.updated_at = datetime()
        SET n += properties
        WITH n, entity
        WHERE entity.identity = false
        MERGE (p:Project {project_id: $project_id})
        MERGE (p)-[r:HAS_ENTITY]->(n)
        RETURN n
    `

	// GET_ENTITIES_BY_PROJECT_TYPE_QUERY fetches entities by project and type
	GET_ENTITIES_BY_PROJECT_TYPE_QUERY = `
        MATCH (p:Project {project_id: $project_id})-[r]-(n:%s)
        RETURN n.name, n.description
    `

	// GET_ENTITIES_BY_TYPE_QUERY fetches all entities of a specific type
	GET_ENTITIES_BY_TYPE_QUERY = `
        MATCH (n:%s)
        RETURN n.name, n.description
    `

	// CREATE_IDENTITY_ENTITY_QUERY creates a new identity-based entity
	CREATE_IDENTITY_ENTITY_QUERY = `
        CREATE (n:%s {
            project_id: $project_id,
            name: $name,
            description: $description
        })
    `

	// MERGE_REUSABLE_ENTITY_QUERY creates or updates a reusable entity
	MERGE_REUSABLE_ENTITY_QUERY = `
        MERGE (n:%s {name: $name})
        ON CREATE SET n.description = $description
        WITH n
        MERGE (p:Project {project_id: $project_id})
        ON MATCH SET p.updated_at = datetime()
        MERGE (p)-[r:HAS_%s]->(n)
        SET r.description = $description
    `

	// DELETE_ENTITY_QUERY removes an entity and its relationships
	DELETE_ENTITY_QUERY = `
        MATCH (n:%s {name: $name})
        OPTIONAL MATCH (n)-[r]-()
        DELETE r, n
    `
)
