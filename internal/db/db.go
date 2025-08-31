package db

import (
	"context"
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	goose "github.com/pressly/goose/v3"
)

type DB interface {
	runMigrations(db *sql.DB, migrationsDir string) error
	Insert(ctx context.Context, url string, shortLink string) error
	Select(ctx context.Context, url string) (string, error)
}

var db *sql.DB

func NewConnection() *sql.DB {
	dsn := `host=127.0.0.1 user=postgres
			password=postgres 
			dbname=postgres port=5432 sslmode=disable`
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	runMigrations(db)
	return db
}

func runMigrations(db *sql.DB) error {
	migrationsDir := "./migrations"
	goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	log.Println("Миграции успешно применены")
	return nil
}

func AddNewAlias(ctx context.Context, url string, shortLink string) error {
	sql, args, err := sq.
		Insert("short_links").
		Columns("original_url", "short_code").
		Values(url, shortLink).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_, err = db.ExecContext(ctx, sql, args...)

	if ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}

func FindByURL(ctx context.Context, url string) (string, error) {
	sql, args, err := sq.
		Select("*").
		From("").
		Where(sq.Eq{"url": url}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var shortCode string
	err = db.QueryRowContext(ctx, sql, args...).Scan(&shortCode)
	if err != nil {
		return "", err
	}

	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	return shortCode, nil
}
