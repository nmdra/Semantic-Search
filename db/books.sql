-- name: InsertBook :exec
INSERT INTO books (title, description, embedding)
VALUES ($1, $2, $3);

-- name: SearchBooks :many
SELECT id, title, description, embedding
FROM books
ORDER BY embedding <=> $1
LIMIT 5;