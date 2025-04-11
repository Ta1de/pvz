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

type ReceptionPostgres struct {
	db *pgx.Conn
}

func NewReceptionPostgres(db *pgx.Conn) *ReceptionPostgres {
	return &ReceptionPostgres{db: db}
}

func (r *ReceptionPostgres) CreateReception(ctx context.Context, pvzId uuid.UUID) (model.Reception, error) {
	query := `
		INSERT INTO reception (pvzId, status)
		VALUES ($1, 'in_progress')
		RETURNING id, dateTime, pvzId, status;
	`

	logger.SugaredLogger.Infow("Creating new reception", "pvzId", pvzId)

	var reception model.Reception
	err := r.db.QueryRow(ctx, query, pvzId).Scan(
		&reception.Id,
		&reception.DateTime,
		&reception.PvzId,
		&reception.Status,
	)
	if err != nil {
		logger.SugaredLogger.Errorw("Failed to create reception", "pvzId", pvzId, "error", err)
		return model.Reception{}, fmt.Errorf("error creating reception: %w", err)
	}

	logger.SugaredLogger.Infow("Successfully created reception", "receptionId", reception.Id)

	return reception, nil
}

func (r *ReceptionPostgres) GetInProgressReception(ctx context.Context, pvzId uuid.UUID) (uuid.UUID, error) {
	query := `
		SELECT id FROM reception
		WHERE pvzId = $1 AND status = 'in_progress'
		ORDER BY datetime DESC
		LIMIT 1;
	`

	var receptionId uuid.UUID
	err := r.db.QueryRow(ctx, query, pvzId).Scan(&receptionId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.SugaredLogger.Warnw("No in-progress reception found", "pvzId", pvzId)
			return uuid.Nil, nil
		}
		logger.SugaredLogger.Errorw("Failed to get in-progress reception", "pvzId", pvzId, "error", err)
		return uuid.Nil, fmt.Errorf("get in-progress reception failed: %w", err)
	}

	logger.SugaredLogger.Infow("Found in-progress reception", "pvzId", pvzId, "receptionId", receptionId)
	return receptionId, nil
}

func (r *ReceptionPostgres) CloseReception(ctx context.Context, pvzId uuid.UUID) error {
	query := `
		UPDATE reception
		SET status = 'close'
		WHERE pvzId = $1 AND status = 'in_progress'
	`

	_, err := r.db.Exec(ctx, query, pvzId)
	if err != nil {
		logger.SugaredLogger.Errorw("Failed to close reception", "pvzId", pvzId, "error", err)
		return fmt.Errorf("failed to close reception: %w", err)
	}

	return nil
}
