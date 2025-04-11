package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"pvz/internal/logger"
	"pvz/internal/repository"
	"pvz/internal/repository/model"
)

type ReceptionService struct {
	repoReception repository.Reception
}

func NewReceptionService(repoReception repository.Reception) *ReceptionService {
	return &ReceptionService{repoReception: repoReception}
}

func (s *ReceptionService) CreateReception(ctx context.Context, pvzId uuid.UUID) (model.Reception, error) {
	logger.SugaredLogger.Infow("Checking for existing in-progress reception", "pvzId", pvzId)

	receptionId, err := s.repoReception.GetInProgressReception(ctx, pvzId)
	if err != nil {
		return model.Reception{}, err
	}
	if receptionId != uuid.Nil {
		logger.SugaredLogger.Warnw("Reception already in progress for PVZ", "pvzId", pvzId)
		return model.Reception{}, fmt.Errorf("an in-progress reception already exists for PVZ %s", pvzId)
	}

	logger.SugaredLogger.Infow("Calling repo to create reception", "pvzId", pvzId)

	reception, err := s.repoReception.CreateReception(ctx, pvzId)
	if err != nil {
		logger.SugaredLogger.Errorw("Failed to create reception in service", "pvzId", pvzId, "error", err)
		return model.Reception{}, err
	}

	logger.SugaredLogger.Infow("Successfully created reception", "receptionId", reception.Id)
	return reception, nil
}

func (s *ReceptionService) CloseReception(ctx context.Context, pvzId uuid.UUID) error {
	logger.SugaredLogger.Infow("Attempting to close reception", "pvzId", pvzId)

	receptionId, err := s.repoReception.GetInProgressReception(ctx, pvzId)
	if err != nil {
		logger.SugaredLogger.Errorw("Failed to get in-progress reception", "pvzId", pvzId, "error", err)
		return fmt.Errorf("cannot close reception: reception lookup failed: %w", err)
	}

	if receptionId == uuid.Nil {
		logger.SugaredLogger.Warnw("No active reception found", "pvzId", pvzId)
		return fmt.Errorf("no active reception found for pvz %s", pvzId)
	}

	err = s.repoReception.CloseReception(ctx, pvzId)
	if err != nil {
		logger.SugaredLogger.Errorw("Failed to close reception", "pvzId", pvzId, "error", err)
		return fmt.Errorf("failed to close reception: %w", err)
	}

	logger.SugaredLogger.Infow("Reception closed successfully", "pvzId", pvzId)

	return nil
}
