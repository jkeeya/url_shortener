package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	sql "database/sql"

	_ "github.com/jackc/pgx"
)

var data_file = "data.json"

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

func makeShortLinkForURL(url string) (bool, string) {
	exists, val := checkExistence(url)
	if exists {
		return true, val
	} else {
		short_link := generateShortLink()
		data := map[string]string{
			url: short_link,
		}

		file, _ := os.OpenFile(data_file, os.O_RDWR|os.O_CREATE, 0644)
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.Encode(data)

		return false, short_link
	}
}

// Переименовать
func checkExistence(url string) (bool, string) {
	file, _ := os.OpenFile(data_file, os.O_RDWR|os.O_CREATE, 0644)
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data map[string]string
	_ = decoder.Decode(&data)

	val, ok := data[url]
	if ok {
		return true, val
	} else {
		return false, ""
	}
}

func main() {
	db, err := sql.Open("postgres",
		"postgres:postgres@tcp(127.0.0.1:5432)")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем новый экземпляр роутера
	r := gin.Default()

	// Определяем маршрут для главной страницы
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Привет, Gin!")
	})

	// Запускаем сервер на порту 8080
	r.Run(":8080")
}
