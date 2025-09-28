package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type createAliasRequest struct {
	URL string `json:"url"`
}

type createAliasResponse struct {
	URL       string `json:"url"`
	ShortLink string `json:"short_link"`
	Existed   bool   `json:"existed"`
}

func (h *Handlers) CreateAlias(c echo.Context) error {
	var request createAliasRequest

	if err := c.Bind(&request); err != nil || request.URL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid or missing url"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	existed, alias, err := h.svc.CreateAlias(ctx, request.URL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := createAliasResponse{
		URL:       alias.URL,
		ShortLink: alias.ShortLink,
		Existed:   existed,
	}
	if existed {
		return c.JSON(http.StatusOK, response)
	}
	return c.JSON(http.StatusCreated, response)
}
