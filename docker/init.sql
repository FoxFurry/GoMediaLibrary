DROP DATABASE IF EXISTS medialibrary;
CREATE DATABASE medialibrary;

DROP USER IF EXISTS mediamanager;
CREATE USER 'mediamanager' IDENTITIED BY 'MainManager';
GRANT ALL PRIVILEGES ON DATABASE medialibrary.* TO 'mediamanager';

DROP TABLE IF EXISTS medialibrary.bookstore;
CREATE TABLE medialibrary.bookstore (
    id SERIAL PRIMARY KEY,
    title TEXT,
    author TEXT,
    year INT,
    description TEXT
);

