CREATE TABLE IF NOT EXISTS configuration (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    value FLOAT NOT NULL
);