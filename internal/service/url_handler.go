package service

import (
	"context"
)

type Repo interface {
	FindByURL(ctx context.Context, url string) (string, error)
	FindByShortLink(ctx context.Context, shortLink string) (string, error)
	AddNewAlias(ctx context.Context, url string, shortLink string) error
}

type URLHandler struct {
	repo Repo
}

func NewURLHandler(repo Repo) *URLHandler {
	return &URLHandler{repo: repo}
}

type Alias struct {
	URL       string
	ShortLink string
}

// TODO: если вдруг сгенерированный шорт уже существует для другого url
func (h URLHandler) CreateAlias(ctx context.Context, url string) (bool, Alias, error) {
	exisingShortLink, _ := h.repo.FindByURL(ctx, url)
	if exisingShortLink != "" {
		return true, Alias{URL: url, ShortLink: exisingShortLink}, nil
	} else {
		shortLink := GenerateShortLink()
		_ = h.repo.AddNewAlias(ctx, url, shortLink)
		return false, Alias{URL: url, ShortLink: shortLink}, nil
	}
}

func (h URLHandler) Redirect(ctx context.Context, shortLink string) (string, error) {
	url, err := h.repo.FindByShortLink(ctx, shortLink)
	if err != nil || url == "" {
		return "", err
	} else {
		return url, err
	}
}

func (h URLHandler) GetShortByURL(ctx context.Context, url string) (string, error) {
	short, err := h.repo.FindByURL(ctx, url)
	if err != nil || short == "" {
		return "", err
	} else {
		return short, nil
	}
}
