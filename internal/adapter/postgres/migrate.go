package postgres

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// RunMigrations applies pending database migrations using embedded SQL scripts.
func RunMigrations(databaseURL string, fs embed.FS, dir string) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open sql database for migrations: %w", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			slog.Warn("failed to close migration database connection", "error", closeErr)
		}
	}()

	goose.SetBaseFS(fs)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect for migrations: %w", err)
	}

	slog.Info("running database migrations", "directory", dir)
	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("execute goose migrations up: %w", err)
	}

	slog.Info("database migrations completed successfully")
	return nil
}
