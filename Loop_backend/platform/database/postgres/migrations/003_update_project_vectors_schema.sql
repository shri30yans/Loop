-- Update project_vectors table to use the simplified schema
-- This migration removes metadata columns and keeps only essential vector search fields

-- First, drop the old table if it exists
DROP TABLE IF EXISTS project_vectors;

-- Create the new table with simplified schema
CREATE TABLE project_vectors (
    project_id UUID PRIMARY KEY REFERENCES projects(project_id) ON DELETE CASCADE,
    summary TEXT NOT NULL,
    embedding vector(384) NOT NULL, -- Using 384 dimensions for all-MiniLM-L6-v2 model
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for vector similarity search
CREATE INDEX project_vectors_embedding_idx ON project_vectors USING ivfflat (embedding vector_cosine_ops)
    WITH (lists = 100); -- Adjust lists parameter based on your dataset size

-- Add trigger to update timestamp
CREATE OR REPLACE FUNCTION update_project_vectors_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_project_vectors_timestamp
    BEFORE UPDATE ON project_vectors
    FOR EACH ROW
    EXECUTE FUNCTION update_project_vectors_timestamp();

-- Add comments for documentation
COMMENT ON TABLE project_vectors IS 'Stores project vector embeddings for semantic search';
COMMENT ON COLUMN project_vectors.project_id IS 'Reference to the project';
COMMENT ON COLUMN project_vectors.summary IS 'Project summary used to generate the embedding';
COMMENT ON COLUMN project_vectors.embedding IS 'Vector embedding of the project summary';
