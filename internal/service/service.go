package service

import (
	"context"

	"github.com/google/uuid"
	"pvz/internal/repository"
	"pvz/internal/repository/model"
)

type User interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	LoginUser(ctx context.Context, email, password string) (string, error)
	DummyLogin(ctx context.Context, role string) (string, error)
}

type Pvz interface {
	CreatePvz(ctx context.Context, pvz model.Pvz) (model.Pvz, error)
}

type Reception interface {
	CreateReception(ctx context.Context, pvzId uuid.UUID) (model.Reception, error)
	CloseReception(ctx context.Context, pvzId uuid.UUID) error
}

type Product interface {
	AddProduct(ctx context.Context, pvzId uuid.UUID, productType string) (model.Product, error)
	DeleteLastProduct(ctx context.Context, pvzId uuid.UUID) error
}

type Service struct {
	User
	Pvz
	Reception
	Product
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:      NewUserService(repos.User),
		Pvz:       NewPvzService(repos.Pvz),
		Reception: NewReceptionService(repos.Reception),
		Product:   NewProductService(repos.Product, repos.Reception),
	}
}
