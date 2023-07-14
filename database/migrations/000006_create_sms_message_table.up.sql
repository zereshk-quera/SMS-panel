CREATE TABLE sms_messages (
    id SERIAL PRIMARY KEY,
    sender VARCHAR(255) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    schedule TIMESTAMP,
    delivery_report TEXT,
    created_at TIMESTAMP DEFAULT current_timestamp,
    account_id INT NOT NULL REFERENCES accounts (id)
);