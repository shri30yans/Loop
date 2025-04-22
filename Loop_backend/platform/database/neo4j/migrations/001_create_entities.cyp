// Drop all existing constraints and indexes first
// CALL db.constraints() YIELD name
// WITH name
// CALL db.drop.constraint(name) YIELD name AS dropped
// RETURN dropped;

// CALL db.indexes() YIELD name
// WITH name
// CALL db.drop.index(name) YIELD name AS dropped
// RETURN dropped;

// Create initial entity constraints
CREATE CONSTRAINT project_id IF NOT EXISTS FOR (p:Project) REQUIRE p.project_id IS UNIQUE;
CREATE CONSTRAINT technology_name IF NOT EXISTS FOR (t:Technology) REQUIRE (t.name, t.project_id) IS UNIQUE;
CREATE CONSTRAINT feature_name IF NOT EXISTS FOR (f:Feature) REQUIRE (f.name, f.project_id) IS UNIQUE;
CREATE CONSTRAINT tag_name IF NOT EXISTS FOR (t:Tag) REQUIRE (t.name, t.project_id) IS UNIQUE;
CREATE CONSTRAINT category_name IF NOT EXISTS FOR (c:Category) REQUIRE c.name IS UNIQUE;
CREATE CONSTRAINT stakeholder_id IF NOT EXISTS FOR (s:Stakeholder) REQUIRE s.id IS UNIQUE;
CREATE CONSTRAINT person_id IF NOT EXISTS FOR (p:Person) REQUIRE p.id IS UNIQUE;
CREATE CONSTRAINT methodology_name IF NOT EXISTS FOR (m:Methodology) REQUIRE m.name IS UNIQUE;

// Create indexes for better query performance
CREATE INDEX project_created_at IF NOT EXISTS FOR (p:Project) ON (p.created_at);
CREATE INDEX tag_category IF NOT EXISTS FOR (t:Tag) ON (t.category);
CREATE INDEX feature_type IF NOT EXISTS FOR (f:Feature) ON (f.type);
CREATE INDEX technology_type IF NOT EXISTS FOR (t:Technology) ON (t.type);
CREATE INDEX platform_type IF NOT EXISTS FOR (p:Platform) ON (p.type);
CREATE INDEX methodology_type IF NOT EXISTS FOR (m:Methodology) ON (m.type);
CREATE INDEX stakeholder_category IF NOT EXISTS FOR (s:Stakeholder) ON (s.category);
