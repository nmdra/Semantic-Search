package embed

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/genai"
)

type GeminiEmbedder struct {
	client *genai.Client
	logger *slog.Logger
}

func NewGeminiEmbedder(ctx context.Context, logger *slog.Logger) (*GeminiEmbedder, error) {

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &GeminiEmbedder{
		client: client,
		logger: logger,
	}, nil
}

func (g *GeminiEmbedder) Embed(ctx context.Context, input string) ([]float32, error) {
	dim := int32(738)
	contents := []*genai.Content{
		genai.NewContentFromText(input, genai.RoleUser),
	}

	resp, err := g.client.Models.EmbedContent(
		ctx,
		"gemini-embedding-001",
		contents,
		&genai.EmbedContentConfig{OutputDimensionality: &dim},
	)
	if err != nil {
		g.logger.Error("embedding failed", "error", err)
		return nil, err
	}

	if len(resp.Embeddings) == 0 {
		return nil, errors.New("no embedding returned")
	}

	embedding := resp.Embeddings[0].Values
	g.logger.Info("Embedding success", "length", len(embedding))
	return embedding, nil
}
