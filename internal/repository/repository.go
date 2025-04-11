package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"pvz/internal/repository/model"
)

type User interface {
	CreateUser(ctx context.Context, user model.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type Pvz interface {
	CreatePvz(ctx context.Context, city string) (model.Pvz, error)
}

type Reception interface {
	CreateReception(ctx context.Context, pvzId uuid.UUID) (model.Reception, error)
	GetInProgressReception(ctx context.Context, pvzId uuid.UUID) (uuid.UUID, error)
	CloseReception(ctx context.Context, pvzId uuid.UUID) error
}

type Product interface {
	CreateProduct(ctx context.Context, product model.Product) (model.Product, error)
	GetLastProductIdByReception(ctx context.Context, receptionId uuid.UUID) (uuid.UUID, error)
	DeleteProductById(ctx context.Context, productId uuid.UUID) error
}

type Repository struct {
	User
	Pvz
	Reception
	Product
}

func NewRepositore(db *pgx.Conn) *Repository {
	return &Repository{
		User:      NewUserPostgres(db),
		Pvz:       NewPvzPostgres(db),
		Reception: NewReceptionPostgres(db),
		Product:   NewProductPostgres(db),
	}
}
