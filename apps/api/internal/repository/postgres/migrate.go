package postgres

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrations embed.FS
var migrationDriver = "pgx"

// RunMigrations applies all pending database migrations.
func RunMigrations(dsn string, logger *zap.Logger) error {
	db, err := sql.Open(migrationDriver, dsn)
	if err != nil {
		return fmt.Errorf("open db for migrations: %w", err)
	}
	defer func() { _ = db.Close() }()

	goose.SetBaseFS(migrations)

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	logger.Info("database migrations applied")
	return nil
}
