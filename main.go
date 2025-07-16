package main

import (
	"context"
	"log/slog"
	"os"
	embed "semantic-search/internal"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.Background()

	embedder, err := embed.NewGeminiEmbedder(ctx, logger)
	if err != nil {
		logger.Error("Failed to initialize embedder", "error", err)
		return
	}

	vector, err := embedder.Embed(ctx, "What is the meaning of life?")
	if err != nil {
		logger.Error("Failed to get embedding", "error", err)
		return
	}

	logger.Info("Embedding length", "length", len(vector))
}
