CREATE TABLE IF NOT EXISTS projects (
    project_id UUID PRIMARY KEY,
    owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    introduction TEXT,
    description TEXT,
    status VARCHAR(50) CHECK (status IN ('draft', 'published', 'completed', 'archived')),
    project_sections JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS project_tags (
    tag_description VARCHAR(50),
    project_id UUID REFERENCES projects(project_id) ON DELETE CASCADE,
    confidence FLOAT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (project_id, tag_description)
);
