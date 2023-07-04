CREATE TABLE IF NOT EXISTS user_numbers (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    number_id INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (number_id) REFERENCES sender_numbers(id)
);