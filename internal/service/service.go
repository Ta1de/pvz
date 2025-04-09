package service

import (
	"context"

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

type Product interface {
}

type Service struct {
	User
	Pvz
	Product
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repos.User),
		Pvz:  NewPvzService(repos.Pvz),
	}
}
