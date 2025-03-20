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
<<<<<<< HEAD
    id UUID PRIMARY KEY,
=======
    id VARCHAR(100) PRIMARY KEY,
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
    username VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    bio TEXT,
    location VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

<<<<<<< HEAD
=======

>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
