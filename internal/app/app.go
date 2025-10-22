package app

import (
	"fmt"
	"strings"

	"github.com/jkeeya/url_shortener/internal/repo/repo_json"
	"github.com/jkeeya/url_shortener/internal/repo/repo_postgres"
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
	RepoPostgres RepositoryType = "repo_postgres"
	RepoJson     RepositoryType = "repo_json"
)

func initRepository(repositoryType RepositoryType, dataSource string) service.Repo {
	switch repositoryType {
	case RepoPostgres:
		return repo_postgres.NewPostgresRepo()
	case RepoJson:
		return repo_json.NewJsonRepo(dataSource)
	}
	return nil
}

func (r *RepositoryType) String() string {
	if r == nil || *r == "" {
		return string(RepoJson)
	}
	return string(*r)
}

func (r *RepositoryType) Set(s string) error {
	switch strings.ToLower(s) {
	case "json":
		*r = RepoJson
	case "postgres", "postgresql", "pg":
		*r = RepoPostgres
	default:
		return fmt.Errorf("invalid repo_type: %q (allowed: json, postgres)", s)
	}
	return nil
}

func NewApp(e *echo.Echo, repo RepositoryType, dataSource string) *App {
	repository := initRepository(repo, dataSource)
	svc := service.NewURLHandler(repository)
	urlHandler := transport.NewHTTPHandlers(svc)
	return &App{
		Echo:       e,
		URLHandler: urlHandler,
		Repo:       &repository,
	}
}
