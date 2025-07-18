DROP INDEX IF EXISTS idx_books_embedding;

-- Create HNSW index on the embedding column
CREATE INDEX IF NOT EXISTS idx_books_embedding_hnsw
ON books USING hnsw (embedding vector_cosine_ops)
WITH (
  m = 16,
  ef_construction = 128
);