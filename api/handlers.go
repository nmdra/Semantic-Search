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
	var req AddBookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	err := h.Service.AddBook(c.Request().Context(), req.Title, req.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": "book added"})
}

func (h *BookHandler) SearchBooks(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing query"})
	}

	results, err := h.Service.SearchBooks(c.Request().Context(), query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}

// TODO: add redis cache
