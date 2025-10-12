package app

import (
	"github.com/jkeeya/url_shortener/internal/repo/postgres"
	"github.com/jkeeya/url_shortener/internal/service"
	"github.com/jkeeya/url_shortener/internal/transport"
	"github.com/labstack/echo/v4"
)

type App struct {
	Echo       *echo.Echo
	URLHandler *transport.Handlers
	Repo       *service.Repo
}
type RepositoryType string

const (
	RepoPostgres RepositoryType = "postgres"
	RepoJson     RepositoryType = "json"
)

func initDatabaseRepository(repositoryType RepositoryType) service.Repo {
	if repositoryType == RepoPostgres {
		return postgres.NewPostgresRepo()
	}
	return nil
}

func NewApp(e *echo.Echo, repo RepositoryType) *App {
	repository := initDatabaseRepository(repo)
	svc := service.NewURLHandler(repository)
	urlHandler := transport.NewHTTPHandlers(svc)
	return &App{
		Echo:       e,
		URLHandler: urlHandler,
		Repo:       &repository,
	}
}
