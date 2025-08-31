package service

import (
	"context"
	"math/rand"

	"github.com/jkeeya/url_shortener/internal/db"
)

func createAlias(ctx context.Context, url string) (bool, string) {
	exists, val := isExists(ctx, url)
	if exists {
		return true, val
	} else {
		short_link := generateShortLink()
		db.AddNewAlias(ctx, url, short_link)
		return false, short_link
	}
}

func generateShortLink() string {
	short_link := [7]byte{}

	// Random ASCII letter code
	letters_range := [][2]int{{65, 90}, {97, 122}}
	for i := range 7 {
		r := letters_range[rand.Intn(len(letters_range))]
		letter := rand.Intn(r[1]-r[0]+1) + r[0]
		short_link[i] = byte(letter)
	}

	result := string(short_link[:])
	return result
}
