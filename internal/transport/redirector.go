package transport

import (
	"context"
	"net/http"
	u "net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) Redirect(c echo.Context) error {
	short := c.Param("short_link")
	if short == "" {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	url, err := h.svc.Redirect(ctx, short)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if url == "" {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	// Приведение URL к абсолютному виде
	normalizedURL, perr := u.Parse(url)
	if perr == nil {
		if normalizedURL.Scheme == "" {
			if strings.HasPrefix(url, "//") {
				url = "https:" + url
			} else {
				url = "https://" + url
			}
		}
	}

	return c.Redirect(http.StatusPermanentRedirect, url)
}
