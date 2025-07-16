## Semantic Book Search with Go, pgvector, and Gemini API

A Golang-based API for semantic search over a book dataset using vector embeddings. Books are embedded using the Gemini API and stored in PostgreSQL with `pgvector`, enabling fast, meaningful similarity search via approximate nearest neighbor indexing (IVFFLAT). Built with Echo, `sqlc`,`go-migrate`, and `pgx`.

