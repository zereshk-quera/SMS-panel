CREATE TABLE budget (
    id INT PRIMARY KEY AUTO_INCREMENT,
    account_id INT NOT NULL,
    amount BIGINT,
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);