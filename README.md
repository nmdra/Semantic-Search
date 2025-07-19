# Semantic Book Search with Go, pgvector, and Gemini API

[![Release with GoReleaser](https://github.com/nmdra/Semantic-Search/actions/workflows/release.yaml/badge.svg)](https://github.com/nmdra/Semantic-Search/actions/workflows/release.yaml)
[![golangci-lint](https://github.com/nmdra/Semantic-Search/actions/workflows/golangci-lint.yaml/badge.svg)](https://github.com/nmdra/Semantic-Search/actions/workflows/golangci-lint.yaml)
[![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Docker Image](https://img.shields.io/badge/docker-ghcr.io%2Fnmdra%2Fsemantic--search-blue?logo=docker)](https://ghcr.io/nmdra/semantic-search)

A Golang-based API for semantic search over a book dataset using vector embeddings. Books are embedded using the Gemini API and stored in PostgreSQL with `pgvector`, enabling fast, meaningful similarity search via approximate nearest neighbor indexing (IVFFLAT). Built with Echo, `sqlc`, `go-migrate`, and `pgx`.

```mermaid
C4Context
  title System Context Diagram for Book Semantic Search
  Enterprise_Boundary(b0, "Book Search Platform") {
    Person(user, "Book Searcher", "User searching for books")
    System(api, "Semantic Search API", "Go service for searching books via embeddings and semantic queries")
    SystemDb(db, "PostgreSQL + pgvector", "Stores books and vector embeddings")
    System_Ext(gemini, "Gemini API", "External API for generating text embeddings")

    Rel(user, api, "Uses", "HTTP\n/Books (POST)\n/Search (GET)")
    Rel(api, db, "Reads & writes", "pgx (SQL)")
    Rel(api, gemini, "Requests text embeddings", "Embed API")
    Rel(gemini, api, "Returns embeddings", "Embedding Response")
    Rel(db, api, "Returns Books/Search Results")
  }
  UpdateRelStyle(user, api, $offsetY="-40", $offsetX="30")
  UpdateRelStyle(api, db, $offsetY="50", $offsetX="-35")
  UpdateRelStyle(api, gemini, $offsetY="-30", $offsetX="-40")
  UpdateRelStyle(gemini, api, $offsetY="0", $offsetX="80")
  UpdateRelStyle(db, api, $offsetY="40", $offsetX="50")
  UpdateLayoutConfig($c4ShapeInRow="3", $c4BoundaryInRow="1")
````

> \[!CAUTION]
> This project is intended for learning and demonstration purposes only.
> While it tries to follow best and security practices, it may contain errors or incomplete implementations.

## Features

* **Semantic Search** — Search books by semantic similarity using vector embeddings
* **Gemini API Integration** — Generates high-quality embeddings via Google's Gemini API
* **PostgreSQL + pgvector** — Efficient storage and approximate nearest neighbor search
* **Redis-powered Cache** — Speeds up repeated search queries with vector caching
* **Run Migrations via CLI** — Run `-migrate` to apply database schema changes at startup
* **Multi-Platform Support** — Build and release for Linux, macOS, Windows, amd64, and arm64
* **Docker & GitHub Container Registry** — Easy deployment with multi-arch Docker images
* **Automated Releases** — GitHub Actions + GoReleaser for continuous delivery

## Getting Started

### Prerequisites

* Go 1.24+
* PostgreSQL with `pgvector` extension installed
* Gemini API Key ([Get API Key Here](https://aistudio.google.com/app/apikey))
* Redis (for vector caching)
* Docker (optional, for containerized deployment)

### API Endpoints

* `POST /books` — Add a book with title and description; stores embedding in DB
* `GET /search?q=your+query` — Search books semantically by query text
* `GET /ping` — Health check endpoint

### Setup PostgreSQL

1. Create your database:

Run migrations:

```bash
make migrate-up
```

Or via the binary:

```bash
semantic-search-api -migrate=true
```
### Environment Variables

```
DATABASE_URL=postgres://user:password@localhost:5432/semantic_search?sslmode=disable
GEMINI_API_KEY=your_gemini_api_key_here
REDIS_URL=localhost:6379
```

### Running Locally

```bash
go run ./cmd/main.go
```

API will be available at `http://localhost:8080`.

### Docker

Run Database migrations:

```bash
docker run --rm \
  --network=host \
  ghcr.io/nmdra/semantic-search:latest \
  -apikey="$GEMINI_API_KEY" \
  -db-dsn="$DATABASE_URL" \
  -migrate
```

Run container:

```bash
docker run --rm \
  --network=host \
  ghcr.io/nmdra/semantic-search:latest \
  -apikey="$GEMINI_API_KEY" \
  -db-dsn="$DATABASE_URL" \
  -redis="localhost:6379" \
  -loglevel="warn"
```
