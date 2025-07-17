package main

import (
	"context"
	"log/slog"
	"os"

	"semantic-search/api"
	"semantic-search/internal/db"
	"semantic-search/internal/embed"
	"semantic-search/internal/repository"
	"semantic-search/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/lmittmann/tint"
)

func main() {
	logger := (slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level: slog.LevelInfo,
	})))
	ctx := context.Background()

	dbpool, err := db.NewPool(ctx)
	if err != nil {
		logger.Error("Failed to connect to DB", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	logger.Info("Connected to PostgreSQL")

	embedder, err := embed.NewGeminiEmbedder(ctx, logger)
	if err != nil {
		logger.Error("Failed to initialize embedder", "error", err)
		return
	}

	repo := repository.New(dbpool)
	bookService := &service.BookService{
		Embedder:   embedder,
		Repository: repo,
	}

	bookHandler := &api.BookHandler{
		Service: bookService,
	}

	e := echo.New()

	// Health check
	e.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	e.GET("/search", bookHandler.SearchBooks)
	e.POST("/books", bookHandler.AddBook)

	e.Logger.Fatal(e.Start(":8080"))
}
