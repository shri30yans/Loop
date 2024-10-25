package database

// SQL queries for database operations
const (
	DropAllTables = `
		DROP TABLE IF EXISTS user_event_participation,
		events,
		project_updates,
		feedback,
		projects,
		users
	`
	CreateUsersTable = `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		email VARCHAR(100) UNIQUE,
		password_hash VARCHAR(100),
		location VARCHAR(100),
		bio TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	CreateProjectsTable = `CREATE TABLE IF NOT EXISTS projects (
		project_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		owner_id INTEGER REFERENCES users(id),
		title VARCHAR(200),
		introduction TEXT,
		description TEXT,
		status VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		tags VARCHAR(200)
	)`

	CreateFeedbackTable = `CREATE TABLE IF NOT EXISTS feedback (
		feedback_id SERIAL PRIMARY KEY,
		project_id UUID REFERENCES projects(project_id),
		user_id INTEGER REFERENCES users(id),
		feedback TEXT
	)`

	CreateProjectUpdatesTable = `CREATE TABLE IF NOT EXISTS project_updates (
		update_id SERIAL PRIMARY KEY,
		project_id UUID REFERENCES projects(project_id),
		title VARCHAR(100),
		body TEXT,
		update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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

	CreateSessionsTables = `CREATE TABLE sessions (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		refresh_token VARCHAR(255) UNIQUE NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
)
