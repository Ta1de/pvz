package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"pvz/internal/repository/model"
)

type User interface {
	CreateUser(ctx context.Context, user model.User) (uuid.UUID, error)
	GetUser(ctx context.Context, email, password string) (model.User, error)
}

type Pvz interface {
	CreatePvz(ctx context.Context, city string) (model.Pvz, error)
}

type Product interface {
}

type Repository struct {
	User
	Pvz
	Product
}

func NewRepositore(db *pgx.Conn) *Repository {
	return &Repository{
		User: NewUserPostgres(db),
		Pvz:  NewPvzPostgres(db),
	}
}
