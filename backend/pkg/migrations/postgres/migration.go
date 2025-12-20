package postgres

import (
	"context"
	"database/sql"
	"embed"
	"log"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	// Игнорируем ошибку, но логируем
	if err := goose.SetDialect("postgres"); err != nil {
		log.Printf("Warning: failed to set dialect (continuing anyway): %v", err)
		// Можно продолжить, если это не критичная ошибка
	}

	return goose.UpContext(context.Background(), db, ".")
}
