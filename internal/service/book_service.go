package service

import (
	"context"
	"fmt"
	"math"
	"semantic-search/internal/embed"
	"semantic-search/internal/repository"

	"github.com/pgvector/pgvector-go"
)

type BookService struct {
	Embedder   embed.Embedder
	Repository *repository.Queries
}

type BookWithSimilarity struct {
	ID          int32
	Title       string
	Description string
	Similarity  float64
}

func (s *BookService) AddBook(ctx context.Context, title, desc string) error {
	vector, err := s.Embedder.Embed(ctx, desc)
	if err != nil {
		return fmt.Errorf("embedding failed: %w", err)
	}

	return s.Repository.InsertBook(ctx, repository.InsertBookParams{
		Title:       title,
		Description: desc,
		Embedding:   pgvector.NewVector(vector),
	})
}

func (s *BookService) SearchBooks(ctx context.Context, query string) ([]BookWithSimilarity, error) {
	vector, err := s.Embedder.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}

	books, err := s.Repository.SearchBooks(ctx, pgvector.NewVector(vector))
	if err != nil {
		return nil, fmt.Errorf("db search failed: %w", err)
	}

	var results []BookWithSimilarity
	for _, book := range books {
		sim, _ := cosineSimilarity(vector, book.Embedding.Slice())
		results = append(results, BookWithSimilarity{
			ID:          book.ID,
			Title:       book.Title,
			Description: book.Description,
			Similarity:  sim,
		})
	}
	return results, nil
}

func cosineSimilarity(a, b []float32) (float64, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vectors must be same length")
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}
	if normA == 0 || normB == 0 {
		return 0, fmt.Errorf("zero vector")
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB)), nil
}
