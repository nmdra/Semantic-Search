package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"semantic-search/internal/embed"
	"semantic-search/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
)

type BookService struct {
	Embedder   embed.Embedder
	Repository *repository.Queries
	Logger     *slog.Logger
}

type BookWithSimilarity struct {
	ID          int32
	ISBN        string
	Title       string
	Description string
	Similarity  float64
}

// AddBook embeds the book description and stores it in the database.
func (s *BookService) AddBook(ctx context.Context, isbn, title, desc string) error {

	// Check Book already exists
	_, err := s.Repository.GetBookByISBN(ctx, pgtype.Text{String: isbn, Valid: true})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("failed to check existing ISBN: %w", err)
	}
	if err == nil {
		return fmt.Errorf("book with isbn %s already exists", isbn)
	}

	vector, err := s.Embedder.Embed(ctx, desc)
	if err != nil {
		return fmt.Errorf("embedding failed: %w", err)
	}

	err = s.Repository.InsertBook(ctx, repository.InsertBookParams{
		Isbn:        pgtype.Text{String: isbn, Valid: true},
		Title:       title,
		Description: desc,
		Embedding:   pgvector.NewVector(vector),
	})
	if err != nil {
		s.Logger.Error("Failed to insert book", "isbn", isbn, "title", title, "error", err)
		return fmt.Errorf("failed to insert book: %w", err)
	}

	return nil
}

// SearchBooks embeds the query, performs vector search, and ranks by cosine similarity.
func (s *BookService) SearchBooks(ctx context.Context, query string) ([]BookWithSimilarity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	vector, err := s.Embedder.Embed(ctx, query)
	if err != nil {
		s.Logger.Error("Embedding failed", "query", query, "error", err)
		return nil, fmt.Errorf("embedding failed: %w", err)
	}

	books, err := s.Repository.SearchBooks(ctx, pgvector.NewVector(vector))
	if err != nil {
		s.Logger.Error("DB search failed", "query", query, "error", err)
		return nil, fmt.Errorf("db search failed: %w", err)
	}

	var results []BookWithSimilarity
	for _, book := range books {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		sim, err := cosineSimilarity(vector, book.Embedding.Slice())
		if err != nil {
			s.Logger.Warn("Failed to compute similarity", "bookID", book.ID, "error", err)
			continue
		}

		results = append(results, BookWithSimilarity{
			ID:          book.ID,
			ISBN:        book.Isbn.String,
			Title:       book.Title,
			Description: book.Description,
			Similarity:  sim,
		})
	}

	return results, nil
}

// cosineSimilarity calculates cosine similarity between two float32 vectors.
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
		return 0, fmt.Errorf("zero vector encountered")
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB)), nil
}
