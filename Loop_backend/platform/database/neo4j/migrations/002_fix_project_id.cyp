// Drop existing Project constraint
DROP CONSTRAINT project_id IF EXISTS;

// Migrate data from id to project_id
MATCH (p:Project)
WHERE p.id IS NOT NULL
WITH p
SET p.project_id = p.id
REMOVE p.id;

// Create new constraint on project_id
CREATE CONSTRAINT project_id IF NOT EXISTS FOR (p:Project) REQUIRE p.project_id IS UNIQUE;
