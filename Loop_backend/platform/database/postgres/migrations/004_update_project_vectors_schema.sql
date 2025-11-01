-- -- Drop existing table if it exists
-- DROP TABLE IF EXISTS project_vectors;

-- Create project_vectors table with updated schema
CREATE TABLE IF NOT EXISTS project_vectors (
    project_id UUID PRIMARY KEY REFERENCES projects(project_id) ON DELETE CASCADE,
    summary TEXT NOT NULL,
    embedding vector(384) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for vector similarity search
CREATE INDEX IF NOT EXISTS project_vectors_embedding_idx ON project_vectors USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_project_vectors_updated_at
    BEFORE UPDATE ON project_vectors
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comment explaining the table
COMMENT ON TABLE project_vectors IS 'Stores project summaries and their vector embeddings for similarity search';
