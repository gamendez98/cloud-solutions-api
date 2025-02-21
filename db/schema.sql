
-- Enable the pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE accounts
(
    id            SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    username      TEXT NOT NULL UNIQUE,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);


CREATE TABLE documents
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name       TEXT    NOT NULL,
    text       TEXT,
    file_path  TEXT,
    embedding VECTOR(384),
    account_id INTEGER NOT NULL REFERENCES accounts (id) ON DELETE CASCADE
);



CREATE TABLE chats
(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    messages JSONB,
    account_id INTEGER NOT NULL REFERENCES accounts (id) ON DELETE CASCADE
);

