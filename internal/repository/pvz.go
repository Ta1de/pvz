package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"pvz/internal/logger"
	"pvz/internal/repository/model"
)

type PvzPostgres struct {
	db *pgx.Conn
}

func NewPvzPostgres(db *pgx.Conn) *PvzPostgres {
	return &PvzPostgres{db: db}
}

func (r *PvzPostgres) CreatePvz(ctx context.Context, city string) (model.Pvz, error) {
	var pvz model.Pvz

	query := `
		INSERT INTO pvz (city)
		VALUES ($1)
		RETURNING id, city, registrationDate
	`

	logger.SugaredLogger.Infow("Inserting new PVZ into database", "city", city)

	err := r.db.QueryRow(ctx, query, city).Scan(&pvz.Id, &pvz.City, &pvz.RegistrationDate)
	if err != nil {
		logger.SugaredLogger.Errorw("Failed to insert PVZ", "city", city, "error", err)
		return pvz, fmt.Errorf("error creating PVZ: %w", err)
	}

	logger.SugaredLogger.Infow("Successfully inserted PVZ", "pvz", pvz)
	return pvz, nil
}
