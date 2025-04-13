package user

import (
	"context"

	"github.com/brnocorreia/api-meu-buzufba/internal/common/dto"
	"github.com/brnocorreia/api-meu-buzufba/pkg/fault"
	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
)

type service struct {
	log      logging.Logger
	userRepo Repository
}

func NewService(log logging.Logger, userRepo Repository) Service {
	return &service{
		log:      log,
		userRepo: userRepo,
	}
}

func (s service) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	userRecord, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve user")
	}
	if userRecord == nil {
		return nil, fault.NewNotFound("user not found")
	}

	user := dto.UserResponse{
		ID:          userRecord.ID,
		Name:        userRecord.Name,
		Username:    userRecord.Username,
		Email:       userRecord.Email,
		IsUfba:      userRecord.IsUfba,
		Activated:   userRecord.Activated,
		ActivatedAt: userRecord.ActivatedAt,
		CreatedAt:   userRecord.CreatedAt,
		UpdatedAt:   userRecord.UpdatedAt,
	}

	return &user, nil
}

func (s service) GetUserByID(ctx context.Context, userId string) (*dto.UserResponse, error) {
	userRecord, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve user")
	} else if userRecord == nil {
		return nil, fault.NewNotFound("user not found")
	}

	user := dto.UserResponse{
		ID:          userRecord.ID,
		Name:        userRecord.Name,
		Username:    userRecord.Username,
		Email:       userRecord.Email,
		IsUfba:      userRecord.IsUfba,
		Activated:   userRecord.Activated,
		ActivatedAt: userRecord.ActivatedAt,
		CreatedAt:   userRecord.CreatedAt,
		UpdatedAt:   userRecord.UpdatedAt,
	}

	return &user, nil
}
