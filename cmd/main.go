package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"semantic-search/api"
	"semantic-search/internal/db"
	"semantic-search/internal/embed"
	"semantic-search/internal/repository"
	"semantic-search/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/tint"
)

type config struct {
	port   int
	apiKey string
	db     struct {
		dsn string
	}
}

func loadConfig() config {
	var cfg config

	defaultAPIKey := os.Getenv("GEMINI_API_KEY")
	defaultDSN := os.Getenv("DATABASE_URL")

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.apiKey, "apikey", defaultAPIKey, "Gemini API Key")
	flag.StringVar(&cfg.db.dsn, "db-dsn", defaultDSN, "PostgreSQL DSN")
	flag.Parse()

	return cfg
}

func main() {
	cfg := loadConfig()

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level: slog.LevelInfo,
	}))

	if cfg.apiKey == "" {
		logger.Error("Gemini API key is required")
		os.Exit(1)
	}

	ctx := context.Background()

	// Database
	dbpool, err := db.NewPool(ctx, cfg.db.dsn, logger)
	if err != nil {
		logger.Error("Failed to connect to DB", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	logger.Info("Connected to PostgreSQL")

	// Embedding client
	embedder, err := embed.NewGeminiEmbedder(ctx, logger, cfg.apiKey)
	if err != nil {
		logger.Error("Failed to initialize embedder", "error", err)
		os.Exit(1)
	}

	// Services and handlers
	repo := repository.New(dbpool)
	bookService := &service.BookService{
		Embedder:   embedder,
		Repository: repo,
		Logger:     logger,
	}
	bookHandler := &api.BookHandler{
		Service: bookService,
	}

	// Echo setup
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})
	e.GET("/search", bookHandler.SearchBooks)
	e.POST("/books", bookHandler.AddBook)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.port)))
}
