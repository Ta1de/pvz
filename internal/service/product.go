package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"pvz/internal/logger"
	"pvz/internal/repository"
	"pvz/internal/repository/model"
	"pvz/metrics"
)

type ProductService struct {
	repoProduct   repository.Product
	repoReception repository.Reception
	logger        logger.Logger
}

func NewProductService(repoProduct repository.Product, repoReception repository.Reception, log logger.Logger) *ProductService {
	return &ProductService{
		repoProduct:   repoProduct,
		repoReception: repoReception,
		logger:        log,
	}
}

func (s *ProductService) AddProduct(ctx context.Context, pvzId uuid.UUID, productType string) (model.Product, error) {
	s.logger.Infow("Adding product", "pvzId", pvzId, "type", productType)

	receptionId, err := s.repoReception.GetInProgressReception(ctx, pvzId)
	if err != nil {
		s.logger.Warnw("Cannot add product, no open reception", "pvzId", pvzId, "error", err)
		return model.Product{}, fmt.Errorf("no open reception for pvz %s: %w", pvzId, err)
	}

	product := model.Product{
		Type:        productType,
		ReceptionId: receptionId,
	}

	created, err := s.repoProduct.CreateProduct(ctx, product)
	if err != nil {
		s.logger.Errorw("Failed to create product", "product", product, "error", err)
		return model.Product{}, fmt.Errorf("failed to create product: %w", err)
	}
	metrics.ProductsAdded.Inc()

	s.logger.Infow("Product created successfully", "productId", created.Id, "receptionId", created.ReceptionId)
	return created, nil
}

func (s *ProductService) DeleteLastProduct(ctx context.Context, pvzId uuid.UUID) error {
	s.logger.Infow("Attempting to delete last product", "pvzId", pvzId)

	receptionId, err := s.repoReception.GetInProgressReception(ctx, pvzId)
	if err != nil {
		s.logger.Errorw("Failed to get in-progress reception", "pvzId", pvzId, "error", err)
		return fmt.Errorf("cannot delete product: reception lookup failed: %w", err)
	}
	if receptionId == uuid.Nil {
		s.logger.Warnw("No active reception found", "pvzId", pvzId)
		return fmt.Errorf("no active reception found for pvz %s", pvzId)
	}

	lastProductId, err := s.repoProduct.GetLastProductIdByReception(ctx, receptionId)
	if err != nil {
		s.logger.Errorw("Failed to get last product ID", "receptionId", receptionId, "error", err)
		return fmt.Errorf("cannot delete product: failed to get last product: %w", err)
	}
	if lastProductId == uuid.Nil {
		s.logger.Warnw("No products found in current reception", "receptionId", receptionId)
		return fmt.Errorf("no products found for current reception")
	}

	err = s.repoProduct.DeleteProductById(ctx, lastProductId)
	if err != nil {
		s.logger.Errorw("Failed to delete product", "productId", lastProductId, "error", err)
		return fmt.Errorf("failed to delete last product: %w", err)
	}

	s.logger.Infow("Product deleted successfully", "productId", lastProductId, "receptionId", receptionId)
	return nil
}
