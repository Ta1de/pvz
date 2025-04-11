package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"pvz/internal/logger"
	"pvz/internal/repository/model"
)

type ProductPostgres struct {
	db *pgx.Conn
}

func NewProductPostgres(db *pgx.Conn) *ProductPostgres {
	return &ProductPostgres{db: db}
}

func (r *ProductPostgres) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	query := `
		INSERT INTO product (type, receptionId)
		VALUES ($1, $2)
		RETURNING id, datetime, type, receptionId;
	`

	var created model.Product
	err := r.db.QueryRow(ctx, query, product.Type, product.ReceptionId).Scan(
		&created.Id,
		&created.DateTime,
		&created.Type,
		&created.ReceptionId,
	)

	if err != nil {
		logger.SugaredLogger.Errorw("Failed to create product", "product", product, "error", err)
		return model.Product{}, fmt.Errorf("error inserting product: %w", err)
	}

	logger.SugaredLogger.Infow("Successfully created product", "product", created)
	return created, nil
}

func (r *ProductPostgres) GetLastProductIdByReception(ctx context.Context, receptionId uuid.UUID) (uuid.UUID, error) {
	query := `
		SELECT id 
		FROM product 
		WHERE receptionId = $1 
		ORDER BY datetime DESC 
		LIMIT 1
	`

	var id uuid.UUID
	err := r.db.QueryRow(ctx, query, receptionId).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.SugaredLogger.Warnw("No products found for reception", "receptionId", receptionId)
			return uuid.Nil, nil
		}
		logger.SugaredLogger.Errorw("Failed to fetch last product ID", "receptionId", receptionId, "error", err)
		return uuid.Nil, fmt.Errorf("failed to get last product id: %w", err)
	}

	logger.SugaredLogger.Infow("Fetched last product ID", "productId", id, "receptionId", receptionId)
	return id, nil
}

func (r *ProductPostgres) DeleteProductById(ctx context.Context, productId uuid.UUID) error {
	query := `DELETE FROM product WHERE id = $1`

	_, err := r.db.Exec(ctx, query, productId)
	if err != nil {
		logger.SugaredLogger.Errorw("Failed to delete product", "productId", productId, "error", err)
		return fmt.Errorf("failed to delete product: %w", err)
	}

	logger.SugaredLogger.Infow("Product deleted successfully", "productId", productId)
	return nil
}
