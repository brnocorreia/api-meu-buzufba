package pg

import (
	"context"
	"fmt"

	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

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
		logging.Error("failed to connect to database", err,
			zap.String("journey", "pg"))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		logging.Error("failed to ping database", err,
			zap.String("journey", "pg"))
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
	logging.Info("initializing database migrations",
		zap.String("journey", "pg"))

	driver, err := postgres.WithInstance(d.db.DB, &postgres.Config{})
	if err != nil {
		logging.Error("failed to migrate database", err,
			zap.String("journey", "pg"))
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/infra/database/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		logging.Error("failed to migrate database", err,
			zap.String("journey", "pg"))
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			logging.Info("database migrations already up to date, skipping...",
				zap.String("journey", "pg"))
			return nil
		}
		logging.Error("failed to migrate database", err,
			zap.String("journey", "pg"))
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	logging.Info("database migrations completed successfully!",
		zap.String("journey", "pg"))
	return nil
}
