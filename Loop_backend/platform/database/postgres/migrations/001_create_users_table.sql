-- DROP TABLE IF EXISTS user_event_participation,
--       passwords,
--       events,
--       sessions,
--       project_sections,
--       comments,
--       projects,
--       project_tags,
--       users
--    ;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    bio TEXT,
    location VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

