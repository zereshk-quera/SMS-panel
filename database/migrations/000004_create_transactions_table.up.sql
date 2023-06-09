CREATE TABLE transactions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    account_id INT NOT NULL,
    amount BIGINT,
    type VARCHAR(255),
    created_at DATETIME,
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);