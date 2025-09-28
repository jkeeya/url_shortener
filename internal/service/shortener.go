package service

import (
	"math/rand"
)

func GenerateShortLink() string {
	shortLink := [7]byte{}

	// Random ASCII letter code
	lettersRange := [][2]int{{65, 90}, {97, 122}}
	for i := range 7 {
		r := lettersRange[rand.Intn(len(lettersRange))]
		letter := rand.Intn(r[1]-r[0]+1) + r[0]
		shortLink[i] = byte(letter)
	}

	result := string(shortLink[:])
	return result
}
