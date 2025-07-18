package db

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a Redis client and logs the connection status.
func NewRedisClient(addr string, logger *slog.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		PoolSize:     10,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis", "addr", addr, "error", err)
		panic("Failed to connect to Redis: " + err.Error())
	}

	logger.Info("Connected to Redis", "addr", addr)
	return client
}
