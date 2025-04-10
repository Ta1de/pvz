package service

import (
	"context"

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
	logger.SugaredLogger.Infow("Calling repo to create reception", "receptionId", pvzId)

	reception, err := s.repoReception.CreateReception(ctx, pvzId)
	if err != nil {
		logger.SugaredLogger.Errorw("Failed to create reception in service", "receptionId", pvzId, "error", err)
		return model.Reception{}, err
	}

	logger.SugaredLogger.Infow("Successfully created reception", "receptionId", reception.Id)
	return reception, nil
}
