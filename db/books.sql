-- name: InsertBook :exec
INSERT INTO books (isbn, title, description, embedding)
VALUES ($1, $2, $3, $4);

-- name: SearchBooks :many
SELECT id, isbn, title, description, embedding
FROM books
ORDER BY embedding <=> $1
LIMIT 5;

-- name: GetBookByISBN :one
SELECT id, isbn
FROM books
WHERE isbn = $1;

-- name: SearchBooksByText :many
SELECT id, isbn, title, description
FROM books
WHERE tsv @@ plainto_tsquery('english', $1)
ORDER BY ts_rank(tsv, plainto_tsquery('english', $1)) DESC
LIMIT 10;