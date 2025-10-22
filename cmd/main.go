package main

import (
	"flag"

	. "github.com/jkeeya/url_shortener/internal/app"
	"github.com/jkeeya/url_shortener/internal/transport"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port := flag.String("port", ":8080", "порт сервера")
	dataSource := flag.String("data_source", "data.json", "путь к файлу базы")

	var repository RepositoryType = RepoJson
	flag.Var(&repository, "repo_type", "тип хранилища")
	flag.Parse()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	app := NewApp(e, repository, *dataSource)

	e.Static("/front/static", "front/static")
	e.File("/", "front/templates/index.html")
	transport.Route(e, app.URLHandler)
	e.Logger.Fatal(e.Start(*port))

}
