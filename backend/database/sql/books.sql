CREATE TABLE author (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    author_id INTEGER REFERENCES author(id)
);