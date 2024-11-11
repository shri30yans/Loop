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
		owner_id SERIAL REFERENCES users(id),
		title VARCHAR(200),
		introduction TEXT,
		description TEXT,
		status VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		tags VARCHAR(200)
	)`

	CreateCommentsTable = `CREATE TABLE IF NOT EXISTS comments (
		comments_id SERIAL PRIMARY KEY,
		project_id SERIAL REFERENCES projects(project_id),
		user_id SERIAL REFERENCES users(id),
		comments TEXT
	)`

	CreateProjectSectionsTable = `CREATE TABLE IF NOT EXISTS project_sections (
    section_id SERIAL,
    project_id INT,
    title VARCHAR(100),
    body TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (section_id, project_id),
	FOREIGN KEY (project_id) REFERENCES projects(project_id)
	)`

	CreateEventsTable = `CREATE TABLE IF NOT EXISTS events (
		event_id SERIAL PRIMARY KEY,
		name VARCHAR(200),
		email VARCHAR(100),
		company VARCHAR(100)
	)`

	CreateUserEventParticipationTable = `CREATE TABLE IF NOT EXISTS user_event_participation (
		user_id SERIAL REFERENCES users(id),
		event_id SERIAL REFERENCES events(event_id),
		PRIMARY KEY (user_id, event_id)
	)`

	CreateSessionsTables = `CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		user_id SERIAL REFERENCES users(id) UNIQUE,
		refresh_token VARCHAR(255) UNIQUE NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
)
