CREATE TABLE accounts (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    budget BIGINT,
    password VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(id)
);