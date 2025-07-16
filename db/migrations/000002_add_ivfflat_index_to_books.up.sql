CREATE INDEX IF NOT EXISTS idx_books_embedding
ON books USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);