CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    budget BIGINT,
    password VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(id)
);