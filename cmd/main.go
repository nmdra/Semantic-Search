package main

import (
	"context"
	"log/slog"
	"os"

	"semantic-search/api"
	embed "semantic-search/internal"
	"semantic-search/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/lmittmann/tint"
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

	embedder, err := embed.NewGeminiEmbedder(ctx, logger)
	if err != nil {
		logger.Error("Failed to initialize embedder", "error", err)
		return
	}

	repo := repository.New(dbpool)
	apiHandler := &api.API{
		Embedder:   embedder,
		Repository: repo,
	}

	e := echo.New()

	// Health check
	e.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	e.GET("/search", apiHandler.SearchBookHandler)
	e.POST("/books", apiHandler.AddBookHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
