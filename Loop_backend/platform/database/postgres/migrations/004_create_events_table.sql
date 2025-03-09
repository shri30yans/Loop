CREATE TABLE IF NOT EXISTS events (
    event_id SERIAL PRIMARY KEY,
    name VARCHAR(200),
    email VARCHAR(100),
    company VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS user_event_participation (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    event_id INTEGER REFERENCES events(event_id),
    PRIMARY KEY (user_id, event_id)
);