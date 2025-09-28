package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type redirectRequest struct {
	ShortLink string `json:"short_link"`
}
type redirectResponse struct {
	Location string `json:"url"`
	Exists   bool   `json:"exists"`
}

func (h *Handlers) Redirect(c echo.Context) error {
	var request redirectRequest
	if err := c.Bind(&request); err != nil || request.ShortLink == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid or missing url"})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	url, err := h.svc.Redirect(ctx, request.ShortLink)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if url == "" {
		// TODO: перенаправить на главную страницу
		url = ""
	}

	response := redirectResponse{
		Location: url,
		Exists:   true,
	}
	return c.JSON(http.StatusPermanentRedirect, response)
}
