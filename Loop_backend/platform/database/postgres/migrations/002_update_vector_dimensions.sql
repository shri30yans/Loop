-- Drop existing table if exists
DROP TABLE IF EXISTS project_vectors;

-- Create simplified table with 384 dimensions
CREATE TABLE project_vectors (
    project_id UUID PRIMARY KEY REFERENCES projects(project_id) ON DELETE CASCADE,
    summary TEXT NOT NULL,
    embedding vector(384) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for similarity search
CREATE INDEX project_vectors_embedding_idx ON project_vectors USING ivfflat (embedding vector_cosine_ops);
