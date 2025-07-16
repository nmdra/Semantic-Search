package main

import (
	"context"
	"log/slog"
	"os"

	embed "semantic-search/internal"
	"semantic-search/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
	"github.com/pgvector/pgvector-go"
)

func main() {
	logger := (slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelInfo})))
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Error("DATABASE_URL not set")
		return
	}

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Error("Failed to connect to DB", "error", err)
		return
	}
	defer dbpool.Close()
	logger.Info("Connected to database")

	queries := repository.New(dbpool)

	embedder, err := embed.NewGeminiEmbedder(ctx, logger)
	if err != nil {
		logger.Error("Failed to initialize embedder", "error", err)
		return
	}

	text := "The Three-Body Problem is a hard sci-fi novel exploring first contact and cosmic civilizations."
	vector, err := embedder.Embed(ctx, text)
	if err != nil {
		logger.Error("Failed to get embedding", "error", err)
		return
	}
	logger.Info("Got embedding", "length", len(vector))

	err = queries.InsertBook(ctx, repository.InsertBookParams{
		Title:       "The Three-Body Problem",
		Description: text,
		Embedding:   pgvector.NewVector(vector),
	})
	if err != nil {
		logger.Error("Failed to insert book into DB", "error", err)
		return
	}

	logger.Info("Book inserted successfully")
}
