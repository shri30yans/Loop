-- Create Tables
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    hashed_password VARCHAR(100),
    location VARCHAR(100),
    bio TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS projects (
    project_id SERIAL PRIMARY KEY,
    owner_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200),
    introduction TEXT,
    description TEXT,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS project_sections (
    section_number INTEGER,
    project_id INTEGER REFERENCES projects(project_id) ON DELETE CASCADE,
    title VARCHAR(100),
    body TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (section_number, project_id)
);

CREATE TABLE IF NOT EXISTS project_tags (
    tag_description VARCHAR(50),
    project_id INTEGER REFERENCES projects(project_id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, tag_description)
);

CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Get User By ID
SELECT 
    u.id,
    u.name, 
    u.email, 
    u.location, 
    u.bio,
    u.created_at,
    (
        SELECT json_agg(json_build_object(
            'project_id', p.project_id,
            'owner_id', p.owner_id,
            'title', p.title,
            'description', p.description,
            'introduction', p.introduction,
            'status', p.status,
            'created_at', to_char(p.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
            'tags', (
                SELECT COALESCE(
                    array_agg(DISTINCT tag_description),
                    '{}'
                )
                FROM project_tags
                WHERE project_id = p.project_id
            )
        ))
        FROM projects p 
        WHERE p.owner_id = u.id
    ) as projects
FROM users u 
WHERE u.id = $1


-- Create Project Function
CREATE OR REPLACE FUNCTION create_project(
       p_title TEXT,
       p_description TEXT,
       p_introduction TEXT,
       p_owner_id INT,
       p_tags TEXT[],
       p_sections JSONB
    ) RETURNS INT AS $$
    DECLARE
       new_project_id INT;
    BEGIN
       -- Insert the main project and get the project_id
       INSERT INTO projects (title, description, introduction, owner_id)
       VALUES (p_title, p_description, p_introduction, p_owner_id)
       RETURNING project_id INTO new_project_id;
       -- Insert tags associated with the project
       INSERT INTO project_tags (project_id, tag_description)
       SELECT new_project_id, unnest(p_tags);
       -- Insert sections associated with the project using JSONB
       INSERT INTO project_sections (project_id, title, body)
       SELECT
          new_project_id,
          section->>'Title',
          section->>'Body'
       FROM jsonb_array_elements(p_sections) AS section;
       -- Return the new project_id
       RETURN new_project_id;
    END;
$$ LANGUAGE plpgsql;


-- Create Project
SELECT create_project($ 1, $ 2, $ 3, $ 4, $ 5 :: text [], $ 6 :: jsonb) 

-- Log Procedure
CREATE TABLE IF NOT EXISTS daily_status_log 
    log_id SERIAL PRIMARY KEY,
    total_users INTEGER NOT NULL,
    total_projects INTEGER NOT NULL,
    logged_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

CREATE OR REPLACE FUNCTION get_projects_and_users_count()
RETURNS TABLE(
    total_projects INT,
    total_users INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
    (SELECT COUNT(*) FROM projects) AS total_projects,
    (SELECT COUNT(*) FROM users) AS total_users;
END;
$$ LANGUAGE plpgsql

-- Fetch Projects based on Keyword
SELECT
    COUNT(*) OVER() AS total_projects,
    p.project_id,
    p.owner_id,
    p.title,
    p.description,
    p.status,
    p.created_at,
    COALESCE(
        json_agg(
            DISTINCT pt.tag_description
        ) FILTER (
            WHERE
                pt.tag_description IS NOT NULL
        ),
        '[]'
    ) as tags
FROM
    projects p
    LEFT JOIN project_tags pt ON p.project_id = pt.project_id
WHERE
    p.title ILIKE '%' || $ 1 || '%'
GROUP BY
    p.project_id,
    p.owner_id,
    p.title,
    p.description,
    p.status,
    p.created_at


-- Fetch All projects
SELECT
    COUNT(*) OVER() AS total_projects,
    p.project_id,
    p.owner_id,
    p.title,
    p.description,
    p.status,
    p.created_at,
    COALESCE(
        json_agg(
            DISTINCT pt.tag_description
        ) FILTER (
            WHERE
                pt.tag_description IS NOT NULL
        ),
        '[]'
    ) as tags
FROM
    projects p
    LEFT JOIN project_tags pt ON p.project_id = pt.project_id
GROUP BY
    p.project_id,
    p.owner_id,
    p.title,
    p.description,
    p.status,
    p.created_at 

-- Fetch Project Info
SELECT
    p.project_id,
    p.owner_id,
    p.title,
    p.description,
    p.introduction,
    p.status,
    p.created_at,
    COALESCE(u.name, '') as owner_name,
    COALESCE(u.email, '') as owner_email,
    u.bio as owner_bio,
    u.location as owner_location,
    COALESCE(
        (
            SELECT
                json_agg(
                    DISTINCT jsonb_build_object(
                        'section_id',
                        ps2.section_number,
                        'title',
                        ps2.title,
                        'body',
                        ps2.body
                    )
                )
            FROM
                project_sections ps2
            WHERE
                ps2.project_id = p.project_id
                AND ps2.section_number IS NOT NULL
        ),
        '[]' :: json
    ) AS sections,
    COALESCE(
        (
            SELECT
                json_agg(DISTINCT pt2.tag_description)
            FROM
                project_tags pt2
            WHERE
                pt2.project_id = p.project_id
                AND pt2.tag_description IS NOT NULL
        ),
        '[]' :: json
    ) AS tags
FROM
    projects p
    LEFT JOIN users u ON p.owner_id = u.id
WHERE
    p.project_id = $ 1
GROUP BY
    p.project_id,
    p.owner_id,
    p.title,
    p.description,
    p.introduction,
    p.status,
    p.created_at,
    u.name,
    u.email,
    u.bio,
    u.location


-- Create Session
INSERT INTO sessions (user_id, refresh_token, expires_at) 
    VALUES ($1, $2, $3) 
    RETURNING id, created_at

-- Get Session By Refresh Token
SELECT id, user_id, refresh_token, expires_at, created_at 
         FROM sessions 
         WHERE refresh_token = $1

-- Get Session by user id
SELECT id, user_id, refresh_token, expires_at, created_at FROM sessions WHERE user_id = $1


-- GetUserByEmail
SELECT
    id,
    email,
    hashed_password,
    created_at
FROM
    users
WHERE
    email = $ 1 

-- Create User
INSERT INTO
    users (name, email, hashed_password)
VALUES
    ($ 1, $ 2, $ 3) RETURNING id,
    name,
    email,
    hashed_password 


-- Audit Table
CREATE TABLE IF NOT EXISTS project_audit (
      audit_id SERIAL PRIMARY KEY,
      project_id INTEGER,
      action VARCHAR(50),
      old_title VARCHAR(200),
      new_title VARCHAR(200),
      old_description TEXT,
      new_description TEXT,
      changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   )

-- Audit Function
CREATE OR REPLACE FUNCTION audit_project_changes() 
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
   IF TG_OP = 'UPDATE' THEN
      INSERT INTO project_audit (project_id, action, old_title, new_title, old_description, new_description)
      VALUES (OLD.project_id, 'UPDATE', OLD.title, NEW.title, OLD.description, NEW.description);
    ELSIF TG_OP = 'INSERT' THEN
    INSERT INTO project_audit (project_id, action, new_title, new_description)
    VALUES (NEW.project_id, 'INSERT', NEW.title, NEW.description);
   ELSIF TG_OP = 'DELETE' THEN
      INSERT INTO project_audit (project_id, action, old_title, old_description)
      VALUES (OLD.project_id, 'DELETE', OLD.title, OLD.description);
   END IF;

   RETURN NULL;  
END;
$$;

-- Audit Trigger
CREATE TRIGGER project_audit_trigger
   AFTER INSERT OR UPDATE OR DELETE ON projects
   FOR EACH ROW EXECUTE FUNCTION audit_project_changes();

-- Update password
UPDATE users SET hashed_password = $1 WHERE id = $2

-- Delete Account
DELETE FROM users WHERE id = $1

-- User and Projects Count Table
CREATE TABLE IF NOT EXISTS daily_status_log (
      log_id SERIAL PRIMARY KEY,
      total_users INTEGER NOT NULL,
      total_projects INTEGER NOT NULL,
      logged_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   )

CREATE OR REPLACE FUNCTION get_projects_and_users_count()
   RETURNS TABLE(
     total_projects INT,
     total_users INT
   ) AS $$
   BEGIN
     RETURN QUERY
     SELECT 
       (SELECT COUNT(*) FROM projects) AS total_projects,
       (SELECT COUNT(*) FROM users) AS total_users;
   END;
   $$ LANGUAGE plpgsql;