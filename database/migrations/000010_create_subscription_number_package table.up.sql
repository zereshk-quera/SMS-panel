CREATE TABLE IF NOT EXISTS subscription_number_package (
    id SERIAL PRIMARY KEY,
    title VARCHAR(55) NOT NULL UNIQUE
);