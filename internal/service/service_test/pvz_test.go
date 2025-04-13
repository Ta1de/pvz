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

func TestCreatePvz_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockPvzRepository)
	mockLogger := new(mocks.MockLogger)
	pvzService := service.NewPvzService(mockRepo, nil, nil, mockLogger)

	expectedPvz := model.Pvz{
		Id:               uuid.New(),
		City:             "Moscow",
		RegistrationDate: time.Now(),
	}

	mockRepo.On("CreatePvz", mock.Anything, expectedPvz.City).Return(expectedPvz, nil)
	mockLogger.On("Infow", "Calling repository to create PVZ", "city", expectedPvz.City)
	mockLogger.On("Infow", "Service successfully created PVZ", "pvz", expectedPvz)

	// Act
	result, err := pvzService.CreatePvz(context.Background(), expectedPvz)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedPvz, result)
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreatePvz_Error(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockPvzRepository)
	mockLogger := new(mocks.MockLogger)
	pvzService := service.NewPvzService(mockRepo, nil, nil, mockLogger)

	testPvz := model.Pvz{
		City: "Moscow", // Make sure this matches the mock expectation
	}
	expectedError := errors.New("database error")

	mockRepo.On("CreatePvz", mock.Anything, testPvz.City).Return(model.Pvz{}, expectedError)
	mockLogger.On("Infow", "Calling repository to create PVZ", "city", testPvz.City)
	mockLogger.On("Errorw", "Service failed to create PVZ", "city", mock.Anything, "error", expectedError)

	// Act
	result, err := pvzService.CreatePvz(context.Background(), testPvz)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, model.Pvz{}, result)
	assert.Contains(t, err.Error(), "error creating PVZ")
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetPvzList_Success(t *testing.T) {
	// Arrange
	mockPvzRepo := new(mocks.MockPvzRepository)
	mockReceptionRepo := new(mocks.MockReceptionRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	mockLogger := new(mocks.MockLogger)
	pvzService := service.NewPvzService(mockPvzRepo, mockReceptionRepo, mockProductRepo, mockLogger)

	limit := 10
	offset := 0
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()

	pvzID := uuid.New()
	receptionID := uuid.New()

	pvzList := []model.Pvz{
		{
			Id:               pvzID,
			City:             "Moscow",
			RegistrationDate: time.Now(),
		},
	}

	receptions := []model.Reception{
		{
			Id:       receptionID,
			DateTime: time.Now(),
			PvzId:    pvzID,
			Status:   "received",
		},
	}

	products := []model.Product{
		{
			Id:          uuid.New(),
			DateTime:    time.Now(),
			Type:        "package",
			ReceptionId: receptionID,
		},
	}

	mockPvzRepo.On("GetPvzListByReceptionDate", mock.Anything, limit, offset, &startDate, &endDate).Return(pvzList, nil)
	mockReceptionRepo.On("GetReceptionsByPvzID", mock.Anything, pvzID).Return(receptions, nil)
	mockProductRepo.On("GetProductsByReceptionID", mock.Anything, receptionID).Return(products, nil)

	// Logger expectations
	mockLogger.On("Infow", "Getting Pvz list by reception date",
		"limit", limit, "offset", offset, "startDate", &startDate, "endDate", &endDate)
	mockLogger.On("Infow", "Processing Pvz", "pvzId", pvzID)
	mockLogger.On("Infow", "Processing Reception", "receptionId", receptionID)
	mockLogger.On("Infow", "Completed processing Pvz", "pvzId", pvzID)
	mockLogger.On("Infow", "Successfully retrieved Pvz list", "count", 1)

	// Act
	result, err := pvzService.GetPvzList(context.Background(), limit, offset, &startDate, &endDate)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, pvzID.String(), result[0].Pvz.Id)
	assert.Len(t, result[0].Receptions, 1)
	assert.Len(t, result[0].Receptions[0].Products, 1)

	mockPvzRepo.AssertExpectations(t)
	mockReceptionRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetPvzList_ErrorGettingPvzList(t *testing.T) {
	// Arrange
	mockPvzRepo := new(mocks.MockPvzRepository)
	mockLogger := new(mocks.MockLogger)
	pvzService := service.NewPvzService(mockPvzRepo, nil, nil, mockLogger)

	limit := 10
	offset := 0
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()
	expectedError := errors.New("database error")

	mockPvzRepo.On("GetPvzListByReceptionDate", mock.Anything, limit, offset, &startDate, &endDate).Return(nil, expectedError)
	mockLogger.On("Infow", "Getting Pvz list by reception date",
		"limit", limit, "offset", offset, "startDate", &startDate, "endDate", &endDate)
	mockLogger.On("Errorw", "Failed to get Pvz list", "error", expectedError)

	// Act
	result, err := pvzService.GetPvzList(context.Background(), limit, offset, &startDate, &endDate)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockPvzRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetPvzList_ErrorGettingReceptions(t *testing.T) {
	// Arrange
	mockPvzRepo := new(mocks.MockPvzRepository)
	mockReceptionRepo := new(mocks.MockReceptionRepository)
	mockLogger := new(mocks.MockLogger)
	pvzService := service.NewPvzService(mockPvzRepo, mockReceptionRepo, nil, mockLogger)

	limit := 10
	offset := 0
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()
	expectedError := errors.New("reception error")

	pvzID := uuid.New()
	pvzList := []model.Pvz{
		{
			Id:               pvzID,
			City:             "Moscow",
			RegistrationDate: time.Now(),
		},
	}

	mockPvzRepo.On("GetPvzListByReceptionDate", mock.Anything, limit, offset, &startDate, &endDate).Return(pvzList, nil)
	mockReceptionRepo.On("GetReceptionsByPvzID", mock.Anything, pvzID).Return(nil, expectedError)

	// Logger expectations
	mockLogger.On("Infow", "Getting Pvz list by reception date",
		"limit", limit, "offset", offset, "startDate", &startDate, "endDate", &endDate)
	mockLogger.On("Infow", "Processing Pvz", "pvzId", pvzID)
	mockLogger.On("Errorw", "Failed to get receptions for Pvz", "pvzId", pvzID, "error", expectedError)

	// Act
	result, err := pvzService.GetPvzList(context.Background(), limit, offset, &startDate, &endDate)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockPvzRepo.AssertExpectations(t)
	mockReceptionRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetPvzList_ErrorGettingProducts(t *testing.T) {
	// Arrange
	mockPvzRepo := new(mocks.MockPvzRepository)
	mockReceptionRepo := new(mocks.MockReceptionRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	mockLogger := new(mocks.MockLogger)
	pvzService := service.NewPvzService(mockPvzRepo, mockReceptionRepo, mockProductRepo, mockLogger)

	limit := 10
	offset := 0
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()
	expectedError := errors.New("product error")

	pvzID := uuid.New()
	receptionID := uuid.New()

	pvzList := []model.Pvz{
		{
			Id:               pvzID,
			City:             "Moscow",
			RegistrationDate: time.Now(),
		},
	}

	receptions := []model.Reception{
		{
			Id:       receptionID,
			DateTime: time.Now(),
			PvzId:    pvzID,
			Status:   "received",
		},
	}

	mockPvzRepo.On("GetPvzListByReceptionDate", mock.Anything, limit, offset, &startDate, &endDate).Return(pvzList, nil)
	mockReceptionRepo.On("GetReceptionsByPvzID", mock.Anything, pvzID).Return(receptions, nil)
	mockProductRepo.On("GetProductsByReceptionID", mock.Anything, receptionID).Return(nil, expectedError)

	// Logger expectations
	mockLogger.On("Infow", "Getting Pvz list by reception date",
		"limit", limit, "offset", offset, "startDate", &startDate, "endDate", &endDate)
	mockLogger.On("Infow", "Processing Pvz", "pvzId", pvzID)
	mockLogger.On("Infow", "Processing Reception", "receptionId", receptionID)
	mockLogger.On("Errorw", "Failed to get products for Reception", "receptionId", receptionID, "error", expectedError)

	// Act
	result, err := pvzService.GetPvzList(context.Background(), limit, offset, &startDate, &endDate)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockPvzRepo.AssertExpectations(t)
	mockReceptionRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
