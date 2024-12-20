package database

// SQL queries for database operations
const (
	DropAllTables = `DROP TABLE IF EXISTS user_event_participation,
      events,
      sessions,
      project_sections,
      comments,
      projects,
      project_tags,
      users
   `

	DropProjectTables = `DROP TABLE IF EXISTS project_sections;
      DROP TABLE IF EXISTS comments;
      DROP TABLE IF EXISTS projects;
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
      owner_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
      title VARCHAR(200),
      introduction TEXT,
      description TEXT,
      status VARCHAR(50),
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   )`

	CreateCommentsTable = `CREATE TABLE IF NOT EXISTS comments (
      comments_id SERIAL PRIMARY KEY,
      project_id INTEGER REFERENCES projects(project_id) ON DELETE CASCADE,
      user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
      comments TEXT
   )`

	CreateProjectSectionsTable = `CREATE TABLE IF NOT EXISTS project_sections (
      section_number INTEGER,
      project_id INTEGER,
      title VARCHAR(100),
      body TEXT,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY (section_number, project_id),
      FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE
   )`

	CreateProjectTagsTable = `CREATE TABLE IF NOT EXISTS project_tags (
      tag_description VARCHAR(50),
      project_id INTEGER REFERENCES projects(project_id) ON DELETE CASCADE,
      PRIMARY KEY (project_id, tag_description)
   )`

	CreateEventsTable = `CREATE TABLE IF NOT EXISTS events (
      event_id SERIAL PRIMARY KEY,
      name VARCHAR(200),
      email VARCHAR(100),
      company VARCHAR(100)
   )`

	CreateUserEventParticipationTable = `CREATE TABLE IF NOT EXISTS user_event_participation (
      user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
      event_id INTEGER REFERENCES events(event_id),
      PRIMARY KEY (user_id, event_id)
   )`

	CreateSessionsTables = `CREATE TABLE IF NOT EXISTS sessions (
      id SERIAL PRIMARY KEY,
      user_id INTEGER REFERENCES users(id) ON DELETE CASCADE UNIQUE,
      refresh_token VARCHAR(255) UNIQUE NOT NULL,
      expires_at TIMESTAMP NOT NULL,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   )`

	CreateProjectFunction = `CREATE OR REPLACE FUNCTION create_project(
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
      INSERT INTO project_sections (project_id, title, body, section_number)
      SELECT
        new_project_id,
        section->>'title',
        section->>'body',
        (section->>'section_number')::int
      FROM jsonb_array_elements(p_sections) AS section;

      -- Return the new project_id
      RETURN new_project_id;
   END;
$$ LANGUAGE plpgsql;`

	CreateUserAndProjectsCountTable = `CREATE TABLE IF NOT EXISTS daily_status_log (
      log_id SERIAL PRIMARY KEY,
      total_users INTEGER NOT NULL,
      total_projects INTEGER NOT NULL,
      logged_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   )`

	CreateUsersAndProjectsCountProcedure = `CREATE OR REPLACE FUNCTION get_projects_and_users_count()
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

	CreateAuditTable = `CREATE TABLE IF NOT EXISTS project_audit (
     audit_id SERIAL PRIMARY KEY,
     project_id INTEGER,
     action VARCHAR(50),
     old_title VARCHAR(200),
     new_title VARCHAR(200),
     old_description TEXT,
     new_description TEXT,
     changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   )`

	CreateAuditProcedure = `CREATE OR REPLACE FUNCTION audit_project_changes() 
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
   IF TG_OP = 'UPDATE' THEN
      INSERT INTO project_audit (project_id, action, old_title, new_title, old_description, new_description)
      VALUES (OLD.project_id, 'UPDATE', OLD.title, NEW.title, OLD.description, NEW.description);
   ELSIF TG_OP = 'DELETE' THEN
      INSERT INTO project_audit (project_id, action, old_title, old_description)
      VALUES (OLD.project_id, 'DELETE', OLD.title, OLD.description);
   END IF;

   RETURN NULL;  -- Trigger functions must return something, typically NULL for AFTER triggers
END;
$$;`

	CreateAuditTrigger = `CREATE TRIGGER project_audit_trigger
   AFTER UPDATE OR DELETE ON projects
   FOR EACH ROW EXECUTE FUNCTION audit_project_changes();`
)
