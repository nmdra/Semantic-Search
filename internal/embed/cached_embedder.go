package embed

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/redis/go-redis/v9"
)

type CachedEmbedder struct {
	Base   Embedder
	Redis  *redis.Client
	Logger *slog.Logger
}

func (c *CachedEmbedder) Embed(ctx context.Context, input string) ([]float32, error) {
	normalized := strings.TrimSpace(strings.ToLower(input))
	hash := xxhash.Sum64String(normalized)
	cacheKey := fmt.Sprintf("embed:%x", hash)

	cached, err := c.Redis.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var vec []float32
		if err := json.Unmarshal(cached, &vec); err == nil {
			c.Logger.Info("Embedding cache hit", "query", input)
			return vec, nil
		}
		c.Logger.Warn("Failed to unmarshal cached embedding", "error", err)
	} else if err != redis.Nil {
		// Only log real Redis errors (not cache misses)
		c.Logger.Warn("Redis GET failed", "key", cacheKey, "error", err)
	}

	// Cache Miss
	vec, err := c.Base.Embed(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}

	data, err := json.Marshal(vec)
	if err != nil {
		c.Logger.Warn("Failed to marshal embedding for cache", "error", err)
	} else {
		err := c.Redis.Set(ctx, cacheKey, data, 24*time.Hour).Err()
		if err != nil {
			c.Logger.Warn("Failed to store embedding in Redis", "key", cacheKey, "error", err)
		} else {
			c.Logger.Debug("Cached embedding", "query", input)
		}
	}

	return vec, nil
}
