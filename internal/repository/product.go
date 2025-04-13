package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pvz/internal/logger"
	"pvz/internal/repository/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductPostgres struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewProductPostgres(db *sqlx.DB, log logger.Logger) *ProductPostgres {
	return &ProductPostgres{
		db:     db,
		logger: log,
	}
}

func (r *ProductPostgres) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	query := `
	INSERT INTO product (type, receptionid)
	VALUES ($1, $2)
	RETURNING id, datetime, type, receptionid;
	`
	var created model.Product
	err := r.db.QueryRowxContext(ctx, query, product.Type, product.ReceptionId).StructScan(&created)
	if err != nil {
		r.logger.Errorw("Failed to create product", "product", product, "error", err)
		return model.Product{}, fmt.Errorf("error inserting product: %w", err)
	}

	r.logger.Infow("Successfully created product", "product", created)
	return created, nil
}

func (r *ProductPostgres) GetLastProductIdByReception(ctx context.Context, receptionId uuid.UUID) (uuid.UUID, error) {
	query := `
		SELECT id 
		FROM product 
		WHERE receptionId = $1 
		ORDER BY datetime DESC 
		LIMIT 1;
	`

	var id uuid.UUID
	err := r.db.GetContext(ctx, &id, query, receptionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warnw("No products found for reception", "receptionId", receptionId)
			return uuid.Nil, nil
		}
		r.logger.Errorw("Failed to fetch last product ID", "receptionId", receptionId, "error", err)
		return uuid.Nil, fmt.Errorf("failed to get last product id: %w", err)
	}

	r.logger.Infow("Fetched last product ID", "productId", id, "receptionId", receptionId)
	return id, nil
}

func (r *ProductPostgres) DeleteProductById(ctx context.Context, productId uuid.UUID) error {
	query := `DELETE FROM product WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, productId)
	if err != nil {
		r.logger.Errorw("Failed to delete product", "productId", productId, "error", err)
		return fmt.Errorf("failed to delete product: %w", err)
	}

	r.logger.Infow("Product deleted successfully", "productId", productId)
	return nil
}

func (r *ProductPostgres) GetProductsByReceptionID(ctx context.Context, receptionId uuid.UUID) ([]model.Product, error) {
	query := `SELECT id, datetime, type, receptionId FROM product WHERE receptionId = $1`

	r.logger.Infow("Executing GetProductsByReceptionID query", "receptionId", receptionId)

	var result []model.Product
	err := r.db.SelectContext(ctx, &result, query, receptionId)
	if err != nil {
		r.logger.Errorw("Failed to fetch products", "error", err, "receptionId", receptionId)
		return nil, err
	}

	r.logger.Infow("Successfully retrieved products", "count", len(result), "receptionId", receptionId)
	return result, nil
}
