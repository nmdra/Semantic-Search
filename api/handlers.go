package api

import (
	"fmt"
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

func (a *API) SearchBookHandler(c echo.Context) error {
	query := c.QueryParam("q")
	fmt.Println(query)
	if query == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing query"})
	}

	vector, err := a.Embedder.Embed(c.Request().Context(), query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "embedding field"})
	}

	results, err := a.Repository.SearchBooks(c.Request().Context(), pgvector.NewVector(vector))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "search failed"})
	}

	return c.JSON(http.StatusOK, results)
}

// TODO: add redis cache
