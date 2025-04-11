package session

import (
	"context"

	"github.com/brnocorreia/api-meu-buzufba/internal/common/dto"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/database/model"
)

type Repository interface {
	Insert(ctx context.Context, session model.Session) error
	Update(ctx context.Context, session model.Session) error
	GetByID(ctx context.Context, sessionId string) (*model.Session, error)
	GetAllByUserID(ctx context.Context, userId string) ([]model.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error)
	GetActiveByUserID(ctx context.Context, userId string) (*model.Session, error)
	DeactivateAll(ctx context.Context, userId string) error
	Delete(ctx context.Context, sessionId string) error
}

type Service interface {
	CreateSession(ctx context.Context, input dto.CreateSession) (sessionId string, err error)
	GetAllSessions(ctx context.Context) ([]dto.SessionResponse, error)
	RenewAccessToken(ctx context.Context, refreshToken string) (*dto.RenewAccessToken, error)
}
