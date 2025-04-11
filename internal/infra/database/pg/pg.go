package pg

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	ctx    context.Context
	logger logging.Logger
	db     *sqlx.DB
}

func NewConnection(dsn string) (*Database, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) DB() *sqlx.DB {
	return d.db
}

func (d *Database) Migrate() error {
	d.logger.Info(d.ctx, "🚀 Initializing database migrations 🚀")

	cmd := exec.Command("flyway", "migrate", "-configFiles=flyway.conf")
	resp, err := cmd.CombinedOutput()
	if err != nil {
		d.logger.Error(d.ctx, "🔴 Error while migrating database 🔴")
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	d.logger.Info(d.ctx, string(resp))
	d.logger.Info(d.ctx, "🟢 Database migrations completed successfully! 🟢")
	return nil
}
