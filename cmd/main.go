package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"semantic-search/api"
	"semantic-search/internal/db"
	"semantic-search/internal/embed"
	"semantic-search/internal/repository"
	"semantic-search/internal/service"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/tint"
)

type config struct {
	port     int
	apiKey   string
	migrate  bool
	logLevel string
	db       struct {
		dsn   string
		redis string
	}
}

func main() {
	cfg := loadConfig()
	logger := setupLogger(cfg.logLevel)

	if cfg.migrate {
		err := runMigrations(cfg, logger)
		if err != nil {
			logger.Error("Migration failed", "error", err)
			os.Exit(1)
		}
		return
	}

	if cfg.apiKey == "" {
		logger.Error("Gemini API key is required")
		os.Exit(1)
	}

	ctx := context.Background()

	dbpool, err := db.NewPool(ctx, cfg.db.dsn, logger)
	if err != nil {
		logger.Error("Database connection failed", "error", err)
		os.Exit(1)
	}
	logger.Info("Connected to PostgreSQL")
	defer dbpool.Close()

	embedder, err := newEmbedder(ctx, cfg, logger)
	if err != nil {
		logger.Error("Failed to setup embedder", "error", err)
		os.Exit(1)
	}

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

	// e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.port)))

	quitectx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	srvAddr := fmt.Sprintf(":%d", cfg.port)

	go func() {
		if err := e.Start(srvAddr); err != nil && err != http.ErrServerClosed {
			logger.Error("Unexpected server shutdown", "error", err)
			os.Exit(1)
		}
	}()

	<-quitectx.Done()
	logger.Debug("Interrupt received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server didn't shut down gracefully — something's still hanging!", "error", err)
	} else {
		logger.Info("Server shut down cleanly. Goodbye!")
	}
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
	flag.StringVar(&cfg.logLevel, "loglevel", "info", "Log level (debug|info|warn|error)")

	flag.Parse()

	if cfg.db.dsn == "" {
		fmt.Fprintln(os.Stderr, "Error: --db-dsn is required")
		os.Exit(1)
	}

	if cfg.apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: --apikey is required")
		os.Exit(1)
	}

	return cfg
}

func setupLogger(levelStr string) *slog.Logger {
	level := slog.LevelInfo

	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		fmt.Fprintf(os.Stderr, "invalid log level %q, defaulting to info\n", levelStr)
	}

	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level: level,
	})
	return slog.New(handler)
}

func runMigrations(cfg config, logger *slog.Logger) error {
	err := db.RunMigrations(cfg.db.dsn, logger)
	if err == nil {
		logger.Info("Migrations applied successfully")
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		logger.Error("Postgres error during migration",
			"code", pgErr.Code,
			"message", pgErr.Message,
			"detail", pgErr.Detail,
			"where", pgErr.Where,
		)
	}

	return fmt.Errorf("failed to run migrations: %w", err)
}

func newEmbedder(ctx context.Context, cfg config, logger *slog.Logger) (embed.Embedder, error) {
	base, err := embed.NewGeminiEmbedder(ctx, logger, cfg.apiKey)
	if err != nil {
		return nil, err
	}

	if cfg.db.redis != "" {
		redisClient := db.NewRedisClient(cfg.db.redis, logger)
		return &embed.CachedEmbedder{
			Base:   base,
			Redis:  redisClient,
			Logger: logger,
		}, nil
	}
	return base, nil
}
