package stop

import (
	"context"

	"github.com/brnocorreia/api-meu-buzufba/internal/common/dto"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/database/model"
)

type Repository interface {
	Insert(ctx context.Context, stop model.Stop) error
	Update(ctx context.Context, stop model.Stop) error
	GetByID(ctx context.Context, stopId int) (*model.Stop, error)
	GetBySlug(ctx context.Context, slug string) (*model.Stop, error)
	Inactivate(ctx context.Context, stopId int) error
}

type Service interface {
	GetStopBySlug(ctx context.Context, slug string) (*dto.StopResponse, error)
	GetStopByID(ctx context.Context, stopId int) (*dto.StopResponse, error)
}
