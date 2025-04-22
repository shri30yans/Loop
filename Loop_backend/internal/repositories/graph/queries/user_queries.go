package queries

// Person node operations - these are specific to user functionality
// All other operations should use standard entity queries
const (
    // Core Person operations
    CreatePersonNodeQuery = `
        MERGE (p:Person {id: $userID})
        SET p.username = $username,
            p.role = $role,
            p.created_at = timestamp()
    `

    UpdatePersonProfileQuery = `
        MATCH (p:Person {id: $userID})
        SET p.username = $username,
            p.bio = $bio,
            p.updated_at = timestamp()
    `

    GetPersonProjectsQuery = `
        MATCH (p:Person {id: $userID})-[r:CREATED|CONTRIBUTED_TO]->(project:Project)
        RETURN project.id, project.title, project.description,
               type(r) as relationship_type,
               r.created_at as contribution_date
        ORDER BY contribution_date DESC
    `

    GetPersonContributionsQuery = `
        MATCH (p:Person {id: $userID})-[r]->(entity)
        WHERE type(r) IN ['CREATED', 'CONTRIBUTED_TO', 'HAS_SKILL']
        RETURN type(r) as activity_type,
               entity.name as target_name,
               labels(entity)[0] as target_type,
               r.created_at as timestamp
        ORDER BY r.created_at DESC
        LIMIT $limit
    `

    // Person relationships with other entities
    LinkPersonToEntityQuery = `
        MATCH (p:Person {id: $userID})
        MATCH (e:%s {%s: $entityID})
        MERGE (p)-[r:%s]->(e)
        SET r.created_at = timestamp(),
            r += $properties
    `

    GetPersonRelatedEntitiesQuery = `
        MATCH (p:Person {id: $userID})-[r:%s]->(e:%s)
        RETURN e.name, e.description, properties(r)
    `
)
