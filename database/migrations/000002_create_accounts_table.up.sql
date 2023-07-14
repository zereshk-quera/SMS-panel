CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    budget BIGINT,
    password VARCHAR(255),
    token VARCHAR(1020),
    is_active boolean,
    is_admin boolean,
    FOREIGN KEY (user_id) REFERENCES users(id)
);