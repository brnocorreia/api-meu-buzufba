package stop

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/brnocorreia/api-meu-buzufba/internal/infra/database/model"
	"github.com/brnocorreia/api-meu-buzufba/pkg/fault"
	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &repo{db: db}
}

func (r repo) Insert(ctx context.Context, stop model.Stop) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var query = `
		INSERT INTO stops (
			name,
			slug,
			latitude,
			longitude,
			security_rating,
			is_active,
			created_at,
			updated_at
		) VALUES (
			:name,
			:slug,
			:latitude,
			:longitude,
			:security_rating,
			:is_active,
			:created_at,
			:updated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, stop)
	if err != nil {
		return fault.New("failed to insert stop", fault.WithError(err))
	}

	return nil
}

func (r repo) Update(ctx context.Context, stop model.Stop) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var query = `
		UPDATE stops
		SET
			name = :name,
			slug = :slug,
			latitude = :latitude,
			longitude = :longitude,
			security_rating = :security_rating,
			is_active = :is_active,
			updated_at = :updated_at
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, stop)
	if err != nil {
		return fault.New("failed to update stop", fault.WithError(err))
	}

	return nil
}

func (r repo) GetByID(ctx context.Context, stopId string) (*model.Stop, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var stop model.Stop
	err := r.db.GetContext(ctx, &stop, "SELECT * FROM stops WHERE id = $1 AND is_active = TRUE LIMIT 1", stopId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to get stop by id", fault.WithError(err))
	}

	return &stop, nil
}

func (r repo) GetBySlug(ctx context.Context, slug string) (*model.Stop, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var stop model.Stop
	err := r.db.GetContext(ctx, &stop, "SELECT * FROM stops WHERE slug = $1 AND is_active = TRUE LIMIT 1", slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to get stop by slug", fault.WithError(err))
	}

	return &stop, nil
}

func (r repo) Inactivate(ctx context.Context, stopId string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "UPDATE stops SET is_active = FALSE WHERE id = $1", stopId)
	if err != nil {
		return fault.New("failed to inactivate stop", fault.WithError(err))
	}

	return nil
}
