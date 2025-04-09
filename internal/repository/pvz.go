package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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

	err := r.db.QueryRow(ctx, query, city).Scan(&pvz.Id, &pvz.City, &pvz.RegistrationDate)
	if err != nil {
		return pvz, fmt.Errorf("ошибка при создании ПВЗ: %w", err)
	}

	return pvz, nil
}
