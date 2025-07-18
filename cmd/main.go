package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"semantic-search/api"
	"semantic-search/internal/db"
	"semantic-search/internal/embed"
	"semantic-search/internal/repository"
	"semantic-search/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/tint"
)

type config struct {
	port    int
	apiKey  string
	migrate bool
	db      struct {
		dsn   string
		redis string
	}
}

func main() {
	cfg := loadConfig()
	logger := setupLogger()

	if cfg.migrate {
		runMigrations(cfg, logger)
		return
	}

	if cfg.apiKey == "" {
		logger.Error("Gemini API key is required")
		os.Exit(1)
	}

	ctx := context.Background()
	dbpool := connectToPostgres(ctx, cfg, logger)
	defer dbpool.Close()

	embedder := setupEmbedder(ctx, cfg, logger)
	startServer(cfg, logger, dbpool, embedder)
}

func loadConfig() config {
	var cfg config

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), `
Semantic Search API Server

Usage:
  semantic-search-api [flags]

Flags:
`)
		flag.PrintDefaults()
	}

	defaultAPIKey := os.Getenv("GEMINI_API_KEY")
	defaultDSN := os.Getenv("DATABASE_URL")
	defaultRedisURL := os.Getenv("REDIS_URL")

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.apiKey, "apikey", defaultAPIKey, "Gemini API Key (or set GEMINI_API_KEY env)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", defaultDSN, "PostgreSQL DSN (or set DATABASE_URL env)")
	flag.StringVar(&cfg.db.redis, "redis", defaultRedisURL, "Redis URL (optional, or set REDIS_URL env)")
	flag.BoolVar(&cfg.migrate, "migrate", false, "Run DB migrations and exit")

	flag.Parse()
	return cfg
}

func setupLogger() *slog.Logger {
	return slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level: slog.LevelInfo,
	}))
}

func runMigrations(cfg config, logger *slog.Logger) {
	if err := db.RunMigrations(cfg.db.dsn, logger); err != nil {
		logger.Error("Migration failed", "error", err)
		os.Exit(1)
	}
	logger.Info("Migrations applied successfully")
}

func connectToPostgres(ctx context.Context, cfg config, logger *slog.Logger) *pgxpool.Pool {
	dbpool, err := db.NewPool(ctx, cfg.db.dsn, logger)
	if err != nil {
		logger.Error("Failed to connect to DB", "error", err)
		os.Exit(1)
	}
	logger.Info("Connected to PostgreSQL")
	return dbpool
}

func setupEmbedder(ctx context.Context, cfg config, logger *slog.Logger) embed.Embedder {
	baseEmbedder, err := embed.NewGeminiEmbedder(ctx, logger, cfg.apiKey)
	if err != nil {
		logger.Error("Failed to initialize embedder", "error", err)
		os.Exit(1)
	}

	if cfg.db.redis != "" {
		redisClient := db.NewRedisClient(cfg.db.redis, logger)
		return &embed.CachedEmbedder{
			Base:   baseEmbedder,
			Redis:  redisClient,
			Logger: logger,
		}
	}
	return baseEmbedder
}

func startServer(cfg config, logger *slog.Logger, dbpool *pgxpool.Pool, embedder embed.Embedder) {
	repo := repository.New(dbpool)
	bookService := &service.BookService{
		Embedder:   embedder,
		Repository: repo,
		Logger:     logger,
	}
	bookHandler := &api.BookHandler{
		Service: bookService,
	}

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout:      5 * time.Second,
		ErrorMessage: "Request timed out.",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			logger.Warn("Timeout on route", "path", c.Path(), "error", err)
		},
	}))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})
	e.GET("/search", bookHandler.SearchBooks)
	e.POST("/books", bookHandler.AddBook)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.port)))
}
