package pg

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Database struct {
	ctx context.Context

	db *sqlx.DB
}

func NewConnection(ctx context.Context, dsn string) (*Database, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{ctx: ctx, db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) DB() *sqlx.DB {
	return d.db
}

func (d *Database) Migrate() error {
	slog.Info("🚀 Initializing database migrations 🚀")

	driver, err := postgres.WithInstance(d.db.DB, &postgres.Config{})
	if err != nil {
		slog.Error("🔴 Error while migrating database 🔴")
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/infra/database/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		slog.Error("🔴 Error while migrating database 🔴")
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			slog.Info("🟢 Database migrations already up to date 🟢")
			return nil
		}
		slog.Error("🔴 Error while migrating database 🔴")
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	slog.Info("🟢 Database migrations completed successfully! 🟢")
	return nil
}
