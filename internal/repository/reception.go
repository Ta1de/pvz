package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"pvz/internal/logger"
	"pvz/internal/repository/model"
)

type ReceptionPostgres struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewReceptionPostgres(db *sqlx.DB, log logger.Logger) *ReceptionPostgres {
	return &ReceptionPostgres{
		db:     db,
		logger: log,
	}
}

func (r *ReceptionPostgres) CreateReception(ctx context.Context, pvzId uuid.UUID) (model.Reception, error) {
	query := `
		INSERT INTO reception (pvzId, status)
		VALUES ($1, 'in_progress')
		RETURNING id, dateTime, pvzId, status;
	`

	r.logger.Infow("Creating new reception", "pvzId", pvzId)

	var reception model.Reception
	if err := r.db.QueryRowxContext(ctx, query, pvzId).StructScan(&reception); err != nil {
		r.logger.Errorw("Failed to create reception", "pvzId", pvzId, "error", err)
		return model.Reception{}, fmt.Errorf("error creating reception: %w", err)
	}

	r.logger.Infow("Successfully created reception", "receptionId", reception.Id)
	return reception, nil
}

func (r *ReceptionPostgres) GetInProgressReception(ctx context.Context, pvzId uuid.UUID) (uuid.UUID, error) {
	query := `
		SELECT id FROM reception
		WHERE pvzId = $1 AND status = 'in_progress'
		ORDER BY dateTime DESC
		LIMIT 1;
	`

	var receptionId uuid.UUID
	err := r.db.GetContext(ctx, &receptionId, query, pvzId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warnw("No in-progress reception found", "pvzId", pvzId)
			return uuid.Nil, nil
		}
		r.logger.Errorw("Failed to get in-progress reception", "pvzId", pvzId, "error", err)
		return uuid.Nil, fmt.Errorf("get in-progress reception failed: %w", err)
	}

	r.logger.Infow("Found in-progress reception", "pvzId", pvzId, "receptionId", receptionId)
	return receptionId, nil
}

func (r *ReceptionPostgres) CloseReception(ctx context.Context, pvzId uuid.UUID) error {
	query := `
		UPDATE reception
		SET status = 'close'
		WHERE pvzId = $1 AND status = 'in_progress'
	`

	_, err := r.db.ExecContext(ctx, query, pvzId)
	if err != nil {
		r.logger.Errorw("Failed to close reception", "pvzId", pvzId, "error", err)
		return fmt.Errorf("failed to close reception: %w", err)
	}

	r.logger.Infow("Successfully closed reception(s)", "pvzId", pvzId)
	return nil
}

func (r *ReceptionPostgres) GetReceptionsByPvzID(ctx context.Context, pvzId uuid.UUID) ([]model.Reception, error) {
	query := `SELECT id, dateTime, pvzId, status FROM reception WHERE pvzId = $1`

	r.logger.Infow("Executing GetReceptionsByPvzID query", "pvzId", pvzId)

	var receptions []model.Reception
	err := r.db.SelectContext(ctx, &receptions, query, pvzId)
	if err != nil {
		r.logger.Errorw("Failed to execute query in GetReceptionsByPvzID", "error", err, "pvzId", pvzId)
		return nil, err
	}

	r.logger.Infow("Successfully retrieved receptions list", "count", len(receptions), "pvzId", pvzId)
	return receptions, nil
}
