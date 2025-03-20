CREATE TABLE IF NOT EXISTS passwords (
<<<<<<< HEAD
    user_id UUID REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    hashed_password CHAR(60) NOT NULL
=======
    user_id VARCHAR(100) REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    hashed_password VARCHAR(255) NOT NULL
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
);


-- CREATE TABLE IF NOT EXISTS sessions (
--     user_id VARCHAR(100) REFERENCES users(id) ON DELETE CASCADE UNIQUE,
--     refresh_token VARCHAR(255) UNIQUE NOT NULL,
--     expires_at TIMESTAMP NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );


