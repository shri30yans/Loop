CREATE TABLE IF NOT EXISTS projects (
<<<<<<< HEAD
    project_id UUID PRIMARY KEY,
    owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
=======
    project_id VARCHAR(100) PRIMARY KEY,
    owner_id VARCHAR(100) REFERENCES users(id) ON DELETE CASCADE,
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
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
<<<<<<< HEAD
    project_id UUID REFERENCES projects(project_id) ON DELETE CASCADE,
    confidence FLOAT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (project_id, tag_description)
);
=======
    project_id VARCHAR(100) REFERENCES projects(project_id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, tag_description)
);
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
