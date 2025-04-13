package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"pvz/internal/repository/model"
)

type MockUserPostgres struct {
	mock.Mock
}

func (m *MockUserPostgres) CreateUser(ctx context.Context, user model.User) (uuid.UUID, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockUserPostgres) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(model.User), args.Error(1)
}

type MockPvzRepository struct {
	mock.Mock
}

func (m *MockPvzRepository) CreatePvz(ctx context.Context, city string) (model.Pvz, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(model.Pvz), args.Error(1)
}

func (m *MockPvzRepository) GetPvzListByReceptionDate(ctx context.Context, limit, offset int, startDate, endDate *time.Time) ([]model.Pvz, error) {
	args := m.Called(ctx, limit, offset, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Pvz), args.Error(1)
}

type MockReceptionRepository struct {
	mock.Mock
}

func (m *MockReceptionRepository) GetReceptionsByPvzID(ctx context.Context, pvzId uuid.UUID) ([]model.Reception, error) {
	args := m.Called(ctx, pvzId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Reception), args.Error(1)
}

func (m *MockReceptionRepository) CreateReception(ctx context.Context, pvzId uuid.UUID) (model.Reception, error) {
	args := m.Called(ctx, pvzId)
	return args.Get(0).(model.Reception), args.Error(1)
}

func (m *MockReceptionRepository) GetInProgressReception(ctx context.Context, pvzId uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, pvzId)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockReceptionRepository) CloseReception(ctx context.Context, pvzId uuid.UUID) error {
	args := m.Called(ctx, pvzId)
	return args.Error(0)
}

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetProductsByReceptionID(ctx context.Context, receptionId uuid.UUID) ([]model.Product, error) {
	args := m.Called(ctx, receptionId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *MockProductRepository) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *MockProductRepository) GetLastProductIdByReception(ctx context.Context, receptionId uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, receptionId)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockProductRepository) DeleteProductById(ctx context.Context, productId uuid.UUID) error {
	args := m.Called(ctx, productId)
	return args.Error(0)
}
