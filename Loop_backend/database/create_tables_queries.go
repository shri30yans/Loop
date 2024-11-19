package database

// SQL queries for database operations
const (
	DropAllTables = `
       DROP TABLE IF EXISTS user_event_participation,
       events,
       sessions,
       project_sections,
       comments,
       projects,
       project_tags,
       users CASCADE
    `
	DropProjectTables = `
       DROP TABLE IF EXISTS project_sections CASCADE;
       DROP TABLE IF EXISTS comments CASCADE;
       DROP TABLE IF EXISTS projects CASCADE;
    `
	CreateUsersTable = `CREATE TABLE IF NOT EXISTS users (
       id SERIAL PRIMARY KEY,
       name VARCHAR(100),
       email VARCHAR(100) UNIQUE,
       hashed_password VARCHAR(100),
       location VARCHAR(100),
       bio TEXT,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`

	CreateProjectsTable = `CREATE TABLE IF NOT EXISTS projects (
       project_id SERIAL PRIMARY KEY,
       owner_id INTEGER REFERENCES users(id),
       title VARCHAR(200),
       introduction TEXT,
       description TEXT,
       status VARCHAR(50),
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )` // Removed the trailing comma

	CreateCommentsTable = `CREATE TABLE IF NOT EXISTS comments (
       comments_id SERIAL PRIMARY KEY,
       project_id INTEGER REFERENCES projects(project_id),
       user_id INTEGER REFERENCES users(id),
       comments TEXT
    )`

	CreateProjectSectionsTable = `CREATE TABLE IF NOT EXISTS project_sections (
       section_number INTEGER,
       project_id INTEGER,
       title VARCHAR(100),
       body TEXT,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       PRIMARY KEY (section_number, project_id),
       FOREIGN KEY (project_id) REFERENCES projects(project_id)
    )`

	CreateProjectTagsTable = `CREATE TABLE IF NOT EXISTS project_tags (
       tag_description VARCHAR(50),
       project_id INTEGER REFERENCES projects(project_id),
       PRIMARY KEY (project_id, tag_description)
    )`

	CreateEventsTable = `CREATE TABLE IF NOT EXISTS events (
       event_id SERIAL PRIMARY KEY,
       name VARCHAR(200),
       email VARCHAR(100),
       company VARCHAR(100)
    )`

	CreateUserEventParticipationTable = `CREATE TABLE IF NOT EXISTS user_event_participation (
       user_id INTEGER REFERENCES users(id),
       event_id INTEGER REFERENCES events(event_id),
       PRIMARY KEY (user_id, event_id)
    )`

	CreateSessionsTables = `CREATE TABLE IF NOT EXISTS sessions (
       id SERIAL PRIMARY KEY,
       user_id INTEGER REFERENCES users(id) UNIQUE,
       refresh_token VARCHAR(255) UNIQUE NOT NULL,
       expires_at TIMESTAMP NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`

	CreateProjectSQLFunction = `CREATE OR REPLACE PROCEDURE create_project(
      p_title TEXT,
      p_description TEXT,
      p_introduction TEXT,
      p_owner_id INT,
      p_tags TEXT[],
      p_sections JSONB
   )
   LANGUAGE plpgsql AS $$
   DECLARE
      new_project_id INT;
   BEGIN
      -- Insert the main project and get the project_id
      INSERT INTO projects (title, description, introduction, owner_id)
      VALUES (
         p_title,
         p_description,
         p_introduction,
         p_owner_id
      )
      RETURNING project_id INTO new_project_id;

      -- Insert tags associated with the project
      INSERT INTO project_tags (project_id, tag_description)
      SELECT new_project_id, unnest(p_tags);

      -- Insert sections associated with the project using JSONB
      INSERT INTO project_sections (project_id, title, body, section_number)
      SELECT
         new_project_id,
         section ->> 'title',
         section ->> 'body',
         (section ->> 'section_number')::INT
      FROM jsonb_array_elements(p_sections) AS section;

      -- Optional: Log success (for debugging or monitoring)
      RAISE NOTICE 'Project created with ID %', new_project_id;
   END;
   $$;`

   CreateUsersAndProjectsCountFunction = `
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
   $$ LANGUAGE plpgsql;`
)
