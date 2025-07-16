package api

import (
	"fmt"
	"math"
	"net/http"
	"semantic-search/internal"
	"semantic-search/internal/repository"

	"github.com/labstack/echo/v4"
	"github.com/pgvector/pgvector-go"
)

// TODO: Add Logger
type API struct {
	Embedder   *internal.GeminiEmbedder
	Repository *repository.Queries
}

type AddBookRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (a *API) AddBookHandler(c echo.Context) error {

	req := AddBookRequest{}
	fmt.Println(req)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	vector, err := a.Embedder.Embed(c.Request().Context(), req.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "embedding failed"})
	}

	err = a.Repository.InsertBook(c.Request().Context(), repository.InsertBookParams{
		Title:       req.Title,
		Description: req.Description,
		Embedding:   pgvector.NewVector(vector),
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db insert failed"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": "book added"})
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

type BookWithSimilarity struct {
	ID          int32   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Similarity  float64 `json:"similarity"`
}

func (a *API) SearchBookHandler(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing query"})
	}

	ctx := c.Request().Context()
	queryVector, err := a.Embedder.Embed(ctx, query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "embedding failed"})
	}

	results, err := a.Repository.SearchBooks(ctx, pgvector.NewVector(queryVector))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "search failed"})
	}

	var resp []BookWithSimilarity
	for _, book := range results {
		sim, err := cosineSimilarity(queryVector, book.Embedding.Slice())
		if err != nil {
			sim = 0
		}
		resp = append(resp, BookWithSimilarity{
			ID:          book.ID,
			Title:       book.Title,
			Description: book.Description,
			Similarity:  sim,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

// TODO: add redis cache
