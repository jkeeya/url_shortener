package service

import (
	"context"
)

type DB interface {
	FindByUrl(ctx context.Context, url string) (string, string, error)
	FindByShortLink(ctx context.Context, shortLink string) (string, error)
	AddNewAlias(ctx context.Context, url, alias string) error
}

type URLHandler struct {
	db DB
}

func NewURLHandler(db DB) *URLHandler {
	return &URLHandler{db: db}
}

type Alias struct {
	URL       string
	ShortLink string
}

// TODO: если вдруг сгенерированный шорт уже существует для другого url
func (h URLHandler) CreateAlias(ctx context.Context, url string) (bool, Alias, error) {
	existingURL, exisingShortLink, _ := h.db.FindByUrl(ctx, url)
	if existingURL != "" {
		return true, Alias{URL: existingURL, ShortLink: exisingShortLink}, nil
	} else {
		shortLink := GenerateShortLink()
		h.db.AddNewAlias(ctx, url, shortLink)
		return false, Alias{URL: url, ShortLink: shortLink}, nil
	}
}

func (h URLHandler) Redirect(ctx context.Context, shortLink string) (string, error) {
	url, err := h.db.FindByShortLink(ctx, shortLink)
	if err != nil || url == "" {
		return "", err
	} else {
		return url, err
	}
}
