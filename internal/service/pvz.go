package service

import (
	"context"
	"fmt"
	"time"

	"pvz/internal/api/response"
	"pvz/internal/logger"
	"pvz/internal/repository"
	"pvz/internal/repository/model"
)

type PvzService struct {
	repoPvz       repository.Pvz
	repoReception repository.Reception
	repoProduct   repository.Product
	logger        logger.Logger
}

func NewPvzService(repoPvz repository.Pvz, repoReception repository.Reception,
	repoProduct repository.Product, log logger.Logger) *PvzService {
	return &PvzService{
		repoPvz:       repoPvz,
		repoReception: repoReception,
		repoProduct:   repoProduct,
		logger:        log,
	}
}

func (s *PvzService) CreatePvz(ctx context.Context, pvz model.Pvz) (model.Pvz, error) {
	s.logger.Infow("Calling repository to create PVZ", "city", pvz.City)

	pvz, err := s.repoPvz.CreatePvz(ctx, pvz.City)
	if err != nil {
		s.logger.Errorw("Service failed to create PVZ", "city", pvz.City, "error", err)
		return model.Pvz{}, fmt.Errorf("error creating PVZ: %w", err)
	}

	s.logger.Infow("Service successfully created PVZ", "pvz", pvz)
	return pvz, nil
}

func (s *PvzService) GetPvzList(ctx context.Context, limit, offset int, startDate, endDate *time.Time) ([]response.PvzFullResponse, error) {
	s.logger.Infow("Getting Pvz list by reception date", "limit", limit, "offset", offset, "startDate", startDate, "endDate", endDate)

	pvzList, err := s.repoPvz.GetPvzListByReceptionDate(ctx, limit, offset, startDate, endDate)
	if err != nil {
		s.logger.Errorw("Failed to get Pvz list", "error", err)
		return nil, err
	}

	var fullResponse []response.PvzFullResponse

	for _, pvz := range pvzList {
		s.logger.Infow("Processing Pvz", "pvzId", pvz.Id)

		receptions, err := s.repoReception.GetReceptionsByPvzID(ctx, pvz.Id)
		if err != nil {
			s.logger.Errorw("Failed to get receptions for Pvz", "pvzId", pvz.Id, "error", err)
			return nil, err
		}

		var receptionWrappers []response.ReceptionWrapper

		for _, rec := range receptions {
			s.logger.Infow("Processing Reception", "receptionId", rec.Id)

			products, err := s.repoProduct.GetProductsByReceptionID(ctx, rec.Id)
			if err != nil {
				s.logger.Errorw("Failed to get products for Reception", "receptionId", rec.Id, "error", err)
				return nil, err
			}

			var productResponses []response.ProductResponse
			for _, p := range products {
				productResponses = append(productResponses, response.ProductResponse{
					Id:          p.Id.String(),
					DateTime:    p.DateTime.Format("2006-01-02 15:04:05"),
					Type:        p.Type,
					ReceptionId: p.ReceptionId.String(),
				})
			}

			receptionWrappers = append(receptionWrappers, response.ReceptionWrapper{
				Reception: response.ReceptionResponse{
					Id:       rec.Id.String(),
					DateTime: rec.DateTime.Format("2006-01-02 15:04:05"),
					PvzId:    rec.PvzId.String(),
					Status:   rec.Status,
				},
				Products: productResponses,
			})
		}

		fullResponse = append(fullResponse, response.PvzFullResponse{
			Pvz: response.PvzResponse{
				Id:               pvz.Id.String(),
				RegistrationDate: pvz.RegistrationDate.Format("2006-01-02 15:04:05"),
				City:             pvz.City,
			},
			Receptions: receptionWrappers,
		})

		s.logger.Infow("Completed processing Pvz", "pvzId", pvz.Id)
	}

	s.logger.Infow("Successfully retrieved Pvz list", "count", len(fullResponse))
	return fullResponse, nil
}
