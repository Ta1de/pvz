package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"pvz/internal/logger"
	"pvz/internal/repository/model"

	"github.com/google/uuid"
)

type User interface {
	CreateUser(ctx context.Context, user model.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type Pvz interface {
	CreatePvz(ctx context.Context, city string) (model.Pvz, error)
	GetPvzListByReceptionDate(ctx context.Context, limit, offset int, startDate, endDate *time.Time) ([]model.Pvz, error)
}

type Reception interface {
	CreateReception(ctx context.Context, pvzId uuid.UUID) (model.Reception, error)
	GetInProgressReception(ctx context.Context, pvzId uuid.UUID) (uuid.UUID, error)
	CloseReception(ctx context.Context, pvzId uuid.UUID) error
	GetReceptionsByPvzID(ctx context.Context, pvzId uuid.UUID) ([]model.Reception, error)
}

type Product interface {
	CreateProduct(ctx context.Context, product model.Product) (model.Product, error)
	GetLastProductIdByReception(ctx context.Context, receptionId uuid.UUID) (uuid.UUID, error)
	DeleteProductById(ctx context.Context, productId uuid.UUID) error
	GetProductsByReceptionID(ctx context.Context, receptionId uuid.UUID) ([]model.Product, error)
}

type Repository struct {
	User
	Pvz
	Reception
	Product
}

func NewRepository(db *sqlx.DB, log logger.Logger) *Repository {
	return &Repository{
		User:      NewUserPostgres(db, log),
		Pvz:       NewPvzPostgres(db, log),
		Reception: NewReceptionPostgres(db, log),
		Product:   NewProductPostgres(db, log),
	}
}
