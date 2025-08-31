package service

import (
	"context"

	"github.com/jkeeya/url_shortener/internal/db"
)

type IURLHandler interface {
	checkExistence(url string) (bool, string)
}

type URLHandler struct{}

func isExists(ctx context.Context, url string) (bool, string) {
	val, _ := db.FindByURL(ctx, url)
	if val != "" {
		return true, val
	} else {
		return false, ""
	}
}
