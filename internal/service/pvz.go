package service

import (
	"context"
	"fmt"

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
	pvz, err := s.repoPvz.CreatePvz(ctx, pvz.City)
	if err != nil {
		return model.Pvz{}, fmt.Errorf("ошибка при создании ПВЗ: %w", err)
	}

	return pvz, nil
}
