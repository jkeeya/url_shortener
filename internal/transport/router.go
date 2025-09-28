package transport

import "github.com/jkeeya/url_shortener/internal/service"

type Handlers struct {
	svc *service.URLHandler
}

func NewHTTPHandlers(svc *service.URLHandler) *Handlers {
	return &Handlers{svc: svc}
}

func Route() {

}
