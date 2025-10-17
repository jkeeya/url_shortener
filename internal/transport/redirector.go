package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type redirectRequest struct {
	ShortLink string `repo_json:"short_link"`
}
type redirectResponse struct {
	Location string `repo_json:"url"`
	Exists   bool   `repo_json:"exists"`
}

func (h *Handlers) Redirect(c echo.Context) error {
	var request redirectRequest
	var response redirectResponse
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
		response = redirectResponse{
			Location: "/",
			Exists:   false,
		}
		return c.JSON(http.StatusPermanentRedirect, response)

	}

	response = redirectResponse{
		Location: url,
		Exists:   true,
	}
	return c.JSON(http.StatusPermanentRedirect, response)
}
