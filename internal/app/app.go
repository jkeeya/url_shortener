package app

import (
	"database/sql"

	. "github.com/jkeeya/url_shortener/internal/service"
)

type App struct {
	DB         *sql.DB
	URLHandler URLHandler
}

func NewApp(db *sql.DB) *App {
	urlHandler := URLHandler{}
	return &App{
		DB:         db,
		URLHandler: urlHandler,
	}
}
