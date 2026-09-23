CREATE TABLE urls (
    alias VARCHAR(255) PRIMARY KEY,
    original_url TEXT NOT NULL UNIQUE
);