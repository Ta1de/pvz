package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pvz/internal/repository/model"
	"pvz/internal/service"
	"pvz/mocks"
)

func TestAddProduct_Success(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(mocks.MockReceptionRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	mockLogger := new(mocks.MockLogger)
	productService := service.NewProductService(mockProductRepo, mockReceptionRepo, mockLogger)

	pvzID := uuid.New()
	receptionID := uuid.New()
	productType := "package"
	expectedProduct := model.Product{
		Id:          uuid.New(),
		Type:        productType,
		ReceptionId: receptionID,
		DateTime:    time.Now(),
	}

	mockReceptionRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(receptionID, nil)
	mockProductRepo.On("CreateProduct", mock.Anything, mock.AnythingOfType("model.Product")).Return(expectedProduct, nil)
	mockLogger.On("Infow", "Adding product", "pvzId", pvzID, "type", productType)
	mockLogger.On("Infow", "Product created successfully", "productId", expectedProduct.Id, "receptionId", receptionID)

	// Act
	result, err := productService.AddProduct(context.Background(), pvzID, productType)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, result)
	mockReceptionRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestAddProduct_NoOpenReception(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(mocks.MockReceptionRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	mockLogger := new(mocks.MockLogger)
	productService := service.NewProductService(mockProductRepo, mockReceptionRepo, mockLogger)

	pvzID := uuid.New()
	productType := "package"
	expectedError := errors.New("no reception found")

	mockReceptionRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(uuid.Nil, expectedError)
	mockLogger.On("Infow", "Adding product", "pvzId", pvzID, "type", productType)
	mockLogger.On("Warnw", "Cannot add product, no open reception", "pvzId", pvzID, "error", expectedError)

	// Act
	result, err := productService.AddProduct(context.Background(), pvzID, productType)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, model.Product{}, result)
	assert.Contains(t, err.Error(), "no open reception for pvz")
	mockReceptionRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestAddProduct_CreateProductError(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(mocks.MockReceptionRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	mockLogger := new(mocks.MockLogger)
	productService := service.NewProductService(mockProductRepo, mockReceptionRepo, mockLogger)

	pvzID := uuid.New()
	receptionID := uuid.New()
	productType := "package"
	expectedError := errors.New("create error")

	mockReceptionRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(receptionID, nil)
	mockProductRepo.On("CreateProduct", mock.Anything, mock.AnythingOfType("model.Product")).Return(model.Product{}, expectedError)
	mockLogger.On("Infow", "Adding product", "pvzId", pvzID, "type", productType)
	mockLogger.On("Errorw", "Failed to create product", "product", mock.Anything, "error", expectedError)

	// Act
	result, err := productService.AddProduct(context.Background(), pvzID, productType)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, model.Product{}, result)
	assert.Contains(t, err.Error(), "failed to create product")
	mockReceptionRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteLastProduct_Success(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(mocks.MockReceptionRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	mockLogger := new(mocks.MockLogger)
	productService := service.NewProductService(mockProductRepo, mockReceptionRepo, mockLogger)

	pvzID := uuid.New()
	receptionID := uuid.New()
	lastProductID := uuid.New()

	mockReceptionRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(receptionID, nil)
	mockProductRepo.On("GetLastProductIdByReception", mock.Anything, receptionID).Return(lastProductID, nil)
	mockProductRepo.On("DeleteProductById", mock.Anything, lastProductID).Return(nil)
	mockLogger.On("Infow", "Attempting to delete last product", "pvzId", pvzID)
	mockLogger.On("Infow", "Product deleted successfully", "productId", lastProductID, "receptionId", receptionID)

	// Act
	err := productService.DeleteLastProduct(context.Background(), pvzID)

	// Assert
	assert.NoError(t, err)
	mockReceptionRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteLastProduct_NoActiveReception(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(mocks.MockReceptionRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	mockLogger := new(mocks.MockLogger)
	productService := service.NewProductService(mockProductRepo, mockReceptionRepo, mockLogger)

	pvzID := uuid.New()

	mockReceptionRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(uuid.Nil, nil)
	mockLogger.On("Infow", "Attempting to delete last product", "pvzId", pvzID)
	mockLogger.On("Warnw", "No active reception found", "pvzId", pvzID)

	// Act
	err := productService.DeleteLastProduct(context.Background(), pvzID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active reception found")
	mockReceptionRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteLastProduct_ReceptionLookupError(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(mocks.MockReceptionRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	mockLogger := new(mocks.MockLogger)
	productService := service.NewProductService(mockProductRepo, mockReceptionRepo, mockLogger)

	pvzID := uuid.New()
	expectedError := errors.New("lookup error")

	mockReceptionRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(uuid.Nil, expectedError)
	mockLogger.On("Infow", "Attempting to delete last product", "pvzId", pvzID)
	mockLogger.On("Errorw", "Failed to get in-progress reception", "pvzId", pvzID, "error", expectedError)

	// Act
	err := productService.DeleteLastProduct(context.Background(), pvzID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reception lookup failed")
	mockReceptionRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteLastProduct_NoProductsFound(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(mocks.MockReceptionRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	mockLogger := new(mocks.MockLogger)
	productService := service.NewProductService(mockProductRepo, mockReceptionRepo, mockLogger)

	pvzID := uuid.New()
	receptionID := uuid.New()

	mockReceptionRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(receptionID, nil)
	mockProductRepo.On("GetLastProductIdByReception", mock.Anything, receptionID).Return(uuid.Nil, nil)
	mockLogger.On("Infow", "Attempting to delete last product", "pvzId", pvzID)
	mockLogger.On("Warnw", "No products found in current reception", "receptionId", receptionID)

	// Act
	err := productService.DeleteLastProduct(context.Background(), pvzID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no products found for current reception")
	mockReceptionRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteLastProduct_DeleteError(t *testing.T) {
	// Arrange
	mockReceptionRepo := new(mocks.MockReceptionRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	mockLogger := new(mocks.MockLogger)
	productService := service.NewProductService(mockProductRepo, mockReceptionRepo, mockLogger)

	pvzID := uuid.New()
	receptionID := uuid.New()
	lastProductID := uuid.New()
	expectedError := errors.New("delete error")

	mockReceptionRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(receptionID, nil)
	mockProductRepo.On("GetLastProductIdByReception", mock.Anything, receptionID).Return(lastProductID, nil)
	mockProductRepo.On("DeleteProductById", mock.Anything, lastProductID).Return(expectedError)
	mockLogger.On("Infow", "Attempting to delete last product", "pvzId", pvzID)
	mockLogger.On("Errorw", "Failed to delete product", "productId", lastProductID, "error", expectedError)

	// Act
	err := productService.DeleteLastProduct(context.Background(), pvzID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete last product")
	mockReceptionRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
