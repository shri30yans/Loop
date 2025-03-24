-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Added vector support to project_tags table in the previous migration

-- Create tags table with vector support
CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    description TEXT,
    embedding vector(384), -- Adjust dimension based on Ollama model
    usage_count INT DEFAULT 0,
    confidence FLOAT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tag_relationships (
    tag1_id UUID REFERENCES tags(id),
    tag2_id UUID REFERENCES tags(id),
    strength FLOAT NOT NULL,
    co_occurrences INT DEFAULT 1,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tag1_id, tag2_id)
);

-- Create index for vector similarity search
CREATE INDEX ON tags USING ivfflat (embedding vector_cosine_ops)
    WITH (lists = 100);
