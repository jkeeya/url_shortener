package main

import (
	. "github.com/jkeeya/url_shortener/internal/app"
	"github.com/jkeeya/url_shortener/internal/transport"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	var repository RepositoryType
	repository = RepoPostgres
	app := NewApp(e, repository)

	transport.Route(e, app.URLHandler)
	e.GET("/")

}
