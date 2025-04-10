package repository

import (
	"context"
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
