package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"pvz/internal/api/response"
	"pvz/internal/logger"
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
	GetPvzList(ctx context.Context, limit, offset int, startDate, endDate *time.Time) ([]response.PvzFullResponse, error)
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

func NewService(repos *repository.Repository, log logger.Logger) *Service {
	return &Service{
		User:      NewUserService(repos.User, log),
		Pvz:       NewPvzService(repos.Pvz, repos.Reception, repos.Product, log),
		Reception: NewReceptionService(repos.Reception, log),
		Product:   NewProductService(repos.Product, repos.Reception, log),
	}
}
