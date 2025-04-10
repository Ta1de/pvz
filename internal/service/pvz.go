package service

import (
	"context"
	"fmt"

	"pvz/internal/logger"
	"pvz/internal/repository"
	"pvz/internal/repository/model"
)

type PvzService struct {
	repoPvz repository.Pvz
}

func NewPvzService(repoPvz repository.Pvz) *PvzService {
	return &PvzService{repoPvz: repoPvz}
}

func (s *PvzService) CreatePvz(ctx context.Context, pvz model.Pvz) (model.Pvz, error) {
	logger.SugaredLogger.Infow("Calling repository to create PVZ", "city", pvz.City)

	pvz, err := s.repoPvz.CreatePvz(ctx, pvz.City)
	if err != nil {
		logger.SugaredLogger.Errorw("Service failed to create PVZ", "city", pvz.City, "error", err)
		return model.Pvz{}, fmt.Errorf("error creating PVZ: %w", err)
	}

	logger.SugaredLogger.Infow("Service successfully created PVZ", "pvz", pvz)
	return pvz, nil
}
