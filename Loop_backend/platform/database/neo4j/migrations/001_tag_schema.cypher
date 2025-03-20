// Create constraints for unique identifiers
CREATE CONSTRAINT tag_name IF NOT EXISTS ON (t:Tag) ASSERT t.name IS UNIQUE;
CREATE CONSTRAINT project_id IF NOT EXISTS ON (p:Project) ASSERT p.id IS UNIQUE;
CREATE CONSTRAINT user_id IF NOT EXISTS ON (u:User) ASSERT u.id IS UNIQUE;

// Create indexes for faster lookups
CREATE INDEX tag_category IF NOT EXISTS FOR (t:Tag) ON (t.category);
CREATE INDEX tag_created IF NOT EXISTS FOR (t:Tag) ON (t.created_at);
CREATE INDEX project_created IF NOT EXISTS FOR (p:Project) ON (p.created_at);

// Define relationship types and their properties
// TAG_RELATIONSHIPS:
// - RELATED_TO: Between tags (strength, co_occurrences)
// - USES_TECH: Project to Tag
// - HAS_SKILL: User to Tag (level, years)
// - CREATED: User to Project

// Sample initial categories if needed
CREATE (c:Category {name: 'Frontend'});
CREATE (c:Category {name: 'Backend'});
CREATE (c:Category {name: 'Database'});
CREATE (c:Category {name: 'DevOps'});
CREATE (c:Category {name: 'Mobile'});
CREATE (c:Category {name: 'AI/ML'});
CREATE (c:Category {name: 'Security'});
CREATE (c:Category {name: 'Cloud'});
