CREATE TABLE phone_books (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL,
    name VARCHAR(255),
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);
CREATE TABLE phone_book_numbers (
    id SERIAL PRIMARY KEY,
    phone_book_id INT NOT NULL,
    prefix VARCHAR(255),
    name VARCHAR(255),
    phone VARCHAR(255),
    FOREIGN KEY (phone_book_id) REFERENCES phone_books(id)
);