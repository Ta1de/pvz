package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"pvz/internal/logger"
	"pvz/internal/repository/model"
)

type PvzPostgres struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewPvzPostgres(db *sqlx.DB, log logger.Logger) *PvzPostgres {
	return &PvzPostgres{
		db:     db,
		logger: log,
	}
}

func (r *PvzPostgres) CreatePvz(ctx context.Context, city string) (model.Pvz, error) {
	var pvz model.Pvz

	query := `
		INSERT INTO pvz (city)
		VALUES ($1)
		RETURNING id, city, registrationDate
	`

	r.logger.Infow("Inserting new PVZ into database", "city", city)

	err := r.db.QueryRowxContext(ctx, query, city).StructScan(&pvz)
	if err != nil {
		r.logger.Errorw("Failed to insert PVZ", "city", city, "error", err)
		return pvz, fmt.Errorf("error creating PVZ: %w", err)
	}

	r.logger.Infow("Successfully inserted PVZ", "pvz", pvz)
	return pvz, nil
}

func (r *PvzPostgres) GetPvzListByReceptionDate(ctx context.Context, limit, offset int, startDate, endDate *time.Time) ([]model.Pvz, error) {
	query := `
		SELECT DISTINCT p.id, p.registrationDate, p.city
		FROM pvz p
		JOIN reception r ON r.pvzId = p.id
		WHERE ($1::timestamp IS NULL OR r.dateTime >= $1)
		  AND ($2::timestamp IS NULL OR r.dateTime <= $2)
		ORDER BY p.registrationDate DESC
		LIMIT $3 OFFSET $4
	`

	r.logger.Infow("Executing GetPvzListByReceptionDate query", "startDate", startDate, "endDate", endDate, "limit", limit, "offset", offset)

	var pvzList []model.Pvz
	err := r.db.SelectContext(ctx, &pvzList, query, startDate, endDate, limit, offset)
	if err != nil {
		r.logger.Errorw("Failed to fetch Pvz list", "error", err)
		return nil, err
	}

	r.logger.Infow("Successfully retrieved Pvz list", "count", len(pvzList))
	return pvzList, nil
}
