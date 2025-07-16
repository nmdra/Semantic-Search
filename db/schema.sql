CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    embedding VECTOR(384)
);

CREATE INDEX IF NOT EXISTS idx_books_embedding
ON books USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);