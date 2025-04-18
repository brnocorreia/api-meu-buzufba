package stop

import (
	"context"

	"github.com/brnocorreia/api-meu-buzufba/internal/common/dto"
	"github.com/brnocorreia/api-meu-buzufba/pkg/fault"
	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"go.uber.org/zap"
)

const stopServiceJourney = "stop service"

type ServiceConfig struct {
	StopRepo Repository
}

type service struct {
	stopRepo Repository
}

func NewService(c ServiceConfig) Service {
	return &service{
		stopRepo: c.StopRepo,
	}
}

func (s service) GetStopByID(ctx context.Context, stopId string) (*dto.StopResponse, error) {
	stopRecord, err := s.stopRepo.GetByID(ctx, stopId)
	if err != nil {
		logging.Error("failed to retrieve stop", err,
			zap.String("journey", stopServiceJourney))
		return nil, fault.NewBadRequest("failed to retrieve stop")
	}
	if stopRecord == nil {
		logging.Info("stop not found",
			zap.String("journey", stopServiceJourney),
			zap.String("stop_id", stopId))
		return nil, fault.NewNotFound("stop not found")
	}

	stop := dto.StopResponse{
		ID:             stopRecord.ID,
		Name:           stopRecord.Name,
		Slug:           stopRecord.Slug,
		Latitude:       stopRecord.Latitude,
		Longitude:      stopRecord.Longitude,
		SecurityRating: stopRecord.SecurityRating,
		IsActive:       stopRecord.IsActive,
		CreatedAt:      stopRecord.CreatedAt,
		UpdatedAt:      stopRecord.UpdatedAt,
	}

	return &stop, nil
}

func (s service) GetStopBySlug(ctx context.Context, slug string) (*dto.StopResponse, error) {
	stopRecord, err := s.stopRepo.GetBySlug(ctx, slug)
	if err != nil {
		logging.Error("failed to retrieve stop", err,
			zap.String("journey", stopServiceJourney))
		return nil, fault.NewBadRequest("failed to retrieve stop")
	}
	if stopRecord == nil {
		logging.Info("stop not found",
			zap.String("journey", stopServiceJourney),
			zap.String("slug", slug))
		return nil, fault.NewNotFound("stop not found")
	}

	stop := dto.StopResponse{
		ID:             stopRecord.ID,
		Name:           stopRecord.Name,
		Slug:           stopRecord.Slug,
		Latitude:       stopRecord.Latitude,
		Longitude:      stopRecord.Longitude,
		SecurityRating: stopRecord.SecurityRating,
		IsActive:       stopRecord.IsActive,
		CreatedAt:      stopRecord.CreatedAt,
		UpdatedAt:      stopRecord.UpdatedAt,
	}

	return &stop, nil
}

func (s service) InactivateStop(ctx context.Context, stopId string) error {
	err := s.stopRepo.Inactivate(ctx, stopId)
	if err != nil {
		logging.Error("failed to inactivate stop", err,
			zap.String("journey", stopServiceJourney))
		return fault.NewBadRequest("failed to inactivate stop")
	}
	return nil
}
