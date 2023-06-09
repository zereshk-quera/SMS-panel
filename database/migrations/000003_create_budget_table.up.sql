CREATE TABLE budget (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL,
    amount BIGINT,
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);