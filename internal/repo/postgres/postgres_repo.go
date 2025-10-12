package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jkeeya/url_shortener/internal/service"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func NewPostgresRepo() service.Repo {
	conn, _ := NewConnection()
	return &PostgresDB{conn: conn}
}

type PostgresDB struct {
	conn *sql.DB
}

func NewConnection() (*sql.DB, error) {
	dsn := `host=127.0.0.1 user=postgres
			password=postgres 
			dbname=postgres port=5432 sslmode=disable`
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения базы: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка доступа к базе: %w", err)
	}

	if err := runMigrations(db, "./migrations"); err != nil {
		return nil, fmt.Errorf("ошибка при накате миграций: %w", err)
	}
	return db, nil
}

func runMigrations(db *sql.DB, migrationsDir string) error {
	_ = goose.SetDialect("postgres")

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("ошибка миграции: %w", err)
	}

	log.Println("миграции успешно применены")
	return nil
}

func (d *PostgresDB) AddNewAlias(ctx context.Context, url string, shortLink string) error {
	query, args, err := sq.
		Insert("short_links").
		Columns("original_url", "short_code").
		Values(url, shortLink).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("ошибка построения запроса: %w", err)
	}

	if _, err = d.conn.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("ошибка добавления алиаса: %w", err)
	}

	return nil
}

func (d *PostgresDB) FindByURL(ctx context.Context, url string) (string, error) {
	query, args, err := sq.
		Select("short_code").
		From("short_links").
		Where(sq.Eq{"original_url": url}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("ошибка построения запроса для поиска по url: %w", err)
	}

	var shortCode string
	if err := d.conn.QueryRowContext(ctx, query, args...).Scan(&shortCode); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("ошибка выполнения запроса поиска по url: %w", err)
	}

	return shortCode, nil
}

func (d *PostgresDB) FindByShortLink(ctx context.Context, shortLink string) (string, error) {
	query, args, err := sq.
		Select("original_url").
		From("short_links").
		Where(sq.Eq{"short_code": shortLink}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("ошибка построения запроса для поиска по короткой ссылке: %w", err)
	}

	var originalURL string
	if err := d.conn.QueryRowContext(ctx, query, args...).Scan(&originalURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("ошибка выполнения запроса поиска по короткой ссылке: %w", err)
	}

	return originalURL, nil
}
