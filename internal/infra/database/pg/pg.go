package pg

import (
	"context"
	"fmt"

	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Database struct {
	ctx    context.Context
	logger logging.Logger
	db     *sqlx.DB
}

func NewConnection(ctx context.Context, logger logging.Logger, dsn string) (*Database, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{ctx: ctx, logger: logger, db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) DB() *sqlx.DB {
	return d.db
}

func (d *Database) Migrate() error {
	d.logger.Info(d.ctx, "🚀 Initializing database migrations 🚀")

	driver, err := postgres.WithInstance(d.db.DB, &postgres.Config{})
	if err != nil {
		d.logger.Error(d.ctx, "🔴 Error while migrating database 🔴")
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/infra/database/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		d.logger.Error(d.ctx, "🔴 Error while migrating database 🔴")
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			d.logger.Info(d.ctx, "🟢 Database migrations already up to date 🟢")
			return nil
		}
		d.logger.Error(d.ctx, "🔴 Error while migrating database 🔴")
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	d.logger.Info(d.ctx, "🟢 Database migrations completed successfully! 🟢")
	return nil
}
