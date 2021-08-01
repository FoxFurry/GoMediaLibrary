
DROP TABLE IF EXISTS bookstore;
CREATE TABLE bookstore (
                                        id SERIAL PRIMARY KEY,
                                        title TEXT,
                                        author TEXT,
                                        year INT,
                                        description TEXT
);