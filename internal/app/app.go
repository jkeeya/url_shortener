package app

import (
	"database/sql"
	. "github.com/jkeeya/url_shortener/internal/service"
	"github.com/labstack/echo/v4"
)

type App struct {
	e          *echo.Echo
	DB         *sql.DB
	URLHandler URLHandler
}

func NewApp(db *sql.DB) *App {
	urlHandler := URLHandler{}
	e := echo.New()
	return &App{
		e:          e,
		DB:         db,
		URLHandler: urlHandler,
	}
}
