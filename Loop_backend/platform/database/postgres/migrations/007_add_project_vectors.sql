-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Create project_vectors table
CREATE TABLE IF NOT EXISTS project_vectors (
    project_id UUID PRIMARY KEY REFERENCES projects(project_id) ON DELETE CASCADE,
    summary TEXT,
    technologies TEXT[],
    domains TEXT[],
    stakeholders TEXT[],
    features TEXT[],
    combined_text TEXT,
    embedding vector(1536),  -- Using 1536 dimensions for OpenAI embeddings
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for vector similarity search
CREATE INDEX ON project_vectors USING ivfflat (embedding vector_cosine_ops)
    WITH (lists = 100);  -- Adjust lists parameter based on expected data size

-- Add function to update timestamp on update
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Add trigger to update timestamp
-- CREATE TRIGGER IF NOT EXISTS update_project_vectors_modtime
--     BEFORE UPDATE ON project_vectors
--     FOR EACH ROW
--     EXECUTE FUNCTION update_modified_column();
