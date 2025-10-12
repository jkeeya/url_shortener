package transport

import (
	"github.com/jkeeya/url_shortener/internal/service"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	svc *service.URLHandler
}

func NewHTTPHandlers(svc *service.URLHandler) *Handlers {
	return &Handlers{svc: svc}
}

func Route(e *echo.Echo, h *Handlers) {
	e.POST("/create/:url", h.CreateAlias)
	e.GET("/:short_link", h.Redirect)
}
