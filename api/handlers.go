package api

import (
	"net/http"
	"semantic-search/internal/service"

	"github.com/labstack/echo/v4"
)

type BookHandler struct {
	Service *service.BookService
}

type AddBookRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *BookHandler) AddBook(c echo.Context) error {
	ctx := c.Request().Context()
	var req AddBookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	select {
	case <-ctx.Done():
		return c.JSON(http.StatusRequestTimeout, echo.Map{"error": "request canceled or timed out"})
	default:
	}

	err := h.Service.AddBook(c.Request().Context(), req.Title, req.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": "book added"})
}

func (h *BookHandler) SearchBooks(c echo.Context) error {
	ctx := c.Request().Context()
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing query"})
	}

	select {
	case <-ctx.Done():
		return c.JSON(http.StatusRequestTimeout, echo.Map{"error": "request canceled or timed out"})
	default:
	}

	results, err := h.Service.SearchBooks(ctx, query)
	if err != nil {
		if ctx.Err() != nil {
			return c.JSON(http.StatusRequestTimeout, echo.Map{"error": "request timed out"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}
