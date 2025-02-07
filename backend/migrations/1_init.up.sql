CREATE TABLE IF NOT EXISTS ping_data (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL,
    ip TEXT NOT NULL,
    last_ping TIMESTAMP DEFAULT NOW()
);
