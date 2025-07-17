# Semantic Book Search with Go, pgvector, and Gemini API

[![Release with GoReleaser](https://github.com/nmdra/Semantic-Search/actions/workflows/release.yaml/badge.svg)](https://github.com/nmdra/Semantic-Search/actions/workflows/release.yaml)
[![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Docker Image](https://img.shields.io/badge/docker-ghcr.io%2Fnmdra%2Fsemantic--search-blue?logo=docker)](https://ghcr.io/nmdra/semantic-search)

A Golang-based API for semantic search over a book dataset using vector embeddings. Books are embedded using the Gemini API and stored in PostgreSQL with `pgvector`, enabling fast, meaningful similarity search via approximate nearest neighbor indexing (IVFFLAT). Built with Echo, `sqlc`, `go-migrate`, and `pgx`.

> [!CAUTION]
> This project is intended for learning and demonstration purposes only.
> While it tries to follow best and security practices, it may contain errors or incomplete implementations.

## Features

* **Semantic Search** — Search books by semantic similarity using vector embeddings
* **Gemini API Integration** — Generates high-quality embeddings via Google's Gemini API
* **PostgreSQL + pgvector** — Efficient storage and approximate nearest neighbor search
* **Multi-Platform Support** — Build and release for Linux, macOS, Windows, amd64, and arm64
* **Docker & GitHub Container Registry** — Easy deployment with multi-arch Docker images
* **Clean Architecture** — Modular codebase with separate API, service, repository layers
* **Automated Releases** — GitHub Actions + GoReleaser for continuous delivery

## Architecture & Directory Structure

```
.
├── api                # HTTP handlers (Echo framework)
├── cmd                # Main application entrypoint
├── internal           # Core business logic, embedder, repository implementations
│   ├── embed          # Gemini embedding client
│   └── repository     # Database access layer (sqlc generated)
├── db                 # SQL migration files & schema
├── Dockerfile         # Multi-stage Docker build for scratch image
├── .goreleaser.yml    # Release automation configuration
├── go.mod             # Go modules dependencies
├── Makefile           # Helper commands (build, migrate, test)
└── README.md          # Project documentation (this file)
```

For detailed project layout, see [Go Project Directory Structure](https://gist.github.com/ayoubzulfiqar/9f1a34049332711fddd4d4b2bfd46096).

## Getting Started

### Prerequisites

* Go 1.24+
* PostgreSQL with `pgvector` extension installed
* Gemini API Key ([Get API Key Here](https://aistudio.google.com/app/apikey))
* Docker (optional, for containerized deployment)

### Setup PostgreSQL

1. Create your database:

Run migrations:

```bash
make migrate-up
```

### Environment Variables

Create a `.env` file with the following:

```
DATABASE_URL=postgres://user:password@localhost:5432/semantic_search?sslmode=disable
GEMINI_API_KEY=your_gemini_api_key_here
```

### Running Locally

```bash
go run ./cmd/main.go
```

API will be available at `http://localhost:8080`.

### API Endpoints

* `POST /books` — Add a book with title and description; stores embedding in DB
* `GET /search?q=your+query` — Search books semantically by query text
* `GET /ping` — Health check endpoint

### Docker

Build multi-arch images with GoReleaser or manually:

```bash
docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/nmdra/semantic-search:latest .
```

Run container:

```bash
docker run -p 8080:8080 ghcr.io/nmdra/semantic-search:latest
```
