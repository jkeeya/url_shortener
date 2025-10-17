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
	repository = RepoJson
	app := NewApp(e, repository)

	e.Static("/static", "static")
	e.File("/", "front/templates/index.html")
	transport.Route(e, app.URLHandler)
	e.Logger.Fatal(e.Start(":8080"))

}
