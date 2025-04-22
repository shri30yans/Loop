package queries

// Project Node Operations
const (
    CreateProjectNodeQuery = `
        MERGE (project:Project {project_id: $project_id})
        ON CREATE SET
            project.name = $name,
            project.description = $description,
            project.created_at = datetime()
        ON MATCH SET
            project.name = $name,
            project.description = $description,
            project.updated_at = datetime()
    `

    LinkProjectToOwnerQuery = `
        MERGE (person:Person {id: $owner_id})
        SET person.role = 'owner'
        WITH person
        MERGE (project:Project {project_id: $project_id})
        MERGE (project)-[r:DEVELOPED_BY]->(person)
        SET r.description = 'Project owner',
            r.weight = 10,
            r.category = 'ownership'
    `

    GetProjectNodesQuery = `
        MATCH (p:Project {project_id: $project_id})-[r]-(n)
        RETURN n.name, labels(n)[0], n.description
    `

    GetProjectRelationshipsQuery = `
        MATCH (p:Project {project_id: $project_id})-[r]-(n)-[r2]-(m)
        WHERE n <> p AND m <> p
        RETURN n.name, m.name, type(r2), r2.description, r2.weight, r2.category
    `

    DeleteProjectQuery = `
        MATCH (p:Project {project_id: $project_id})
        DETACH DELETE p
    `
)
