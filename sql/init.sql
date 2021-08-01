CREATE USER bookmanager;

CREATE TABLE bookstore (
    id SERIAL PRIMARY KEY,
    title TEXT,
    author TEXT,
    year INT,
    description TEXT
);

GRANT ALL PRIVILEGES ON DATABASE bookstore TO bookmanaher;