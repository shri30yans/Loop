-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Added vector support to project_tags table in the previous migration

-- Create tags table with vector support
CREATE TABLE IF NOT EXISTS tags (
        id UUID PRIMARY KEY,
        name VARCHAR(50) UNIQUE,
        category VARCHAR(50),
        vector vector(384),
        usage_count INTEGER DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create tag relationships table
CREATE TABLE IF NOT EXISTS tag_relationships (
    tag1_id UUID REFERENCES tags(id),
    tag2_id UUID REFERENCES tags(id),
    strength FLOAT DEFAULT 0.0,
    co_occurrences INTEGER DEFAULT 0,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tag1_id, tag2_id)
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_project_tags_confidence ON project_tags(confidence);
CREATE INDEX IF NOT EXISTS idx_tags_category ON tags(category);
CREATE INDEX IF NOT EXISTS idx_tag_relationships_strength ON tag_relationships(strength);
