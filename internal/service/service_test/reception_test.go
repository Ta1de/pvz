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

func TestCreateReception_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockReceptionRepository)
	mockLogger := new(mocks.MockLogger)
	receptionService := service.NewReceptionService(mockRepo, mockLogger)

	pvzID := uuid.New()
	expectedReception := model.Reception{
		Id:       uuid.New(),
		DateTime: time.Now(),
		PvzId:    pvzID,
		Status:   "in_progress",
	}

	mockRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(uuid.Nil, nil)
	mockRepo.On("CreateReception", mock.Anything, pvzID).Return(expectedReception, nil)
	mockLogger.On("Infow", "Checking for existing in-progress reception", "pvzId", pvzID)
	mockLogger.On("Infow", "Calling repo to create reception", "pvzId", pvzID)
	mockLogger.On("Infow", "Successfully created reception", "receptionId", expectedReception.Id)

	// Act
	result, err := receptionService.CreateReception(context.Background(), pvzID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedReception, result)
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateReception_ExistingReception(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockReceptionRepository)
	mockLogger := new(mocks.MockLogger)
	receptionService := service.NewReceptionService(mockRepo, mockLogger)

	pvzID := uuid.New()
	existingReceptionID := uuid.New()

	mockRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(existingReceptionID, nil)
	mockLogger.On("Infow", "Checking for existing in-progress reception", "pvzId", pvzID)
	mockLogger.On("Warnw", "Reception already in progress for PVZ", "pvzId", pvzID)

	// Act
	result, err := receptionService.CreateReception(context.Background(), pvzID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, model.Reception{}, result)
	assert.Contains(t, err.Error(), "an in-progress reception already exists")
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateReception_GetInProgressError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockReceptionRepository)
	mockLogger := new(mocks.MockLogger)
	receptionService := service.NewReceptionService(mockRepo, mockLogger)

	pvzID := uuid.New()
	expectedError := errors.New("database error")

	mockRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(uuid.Nil, expectedError)
	mockLogger.On("Infow", "Checking for existing in-progress reception", "pvzId", pvzID)

	// Act
	result, err := receptionService.CreateReception(context.Background(), pvzID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, model.Reception{}, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateReception_CreateError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockReceptionRepository)
	mockLogger := new(mocks.MockLogger)
	receptionService := service.NewReceptionService(mockRepo, mockLogger)

	pvzID := uuid.New()
	expectedError := errors.New("create error")

	mockRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(uuid.Nil, nil)
	mockRepo.On("CreateReception", mock.Anything, pvzID).Return(model.Reception{}, expectedError)
	mockLogger.On("Infow", "Checking for existing in-progress reception", "pvzId", pvzID)
	mockLogger.On("Infow", "Calling repo to create reception", "pvzId", pvzID)
	mockLogger.On("Errorw", "Failed to create reception in service", "pvzId", pvzID, "error", expectedError)

	// Act
	result, err := receptionService.CreateReception(context.Background(), pvzID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, model.Reception{}, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCloseReception_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockReceptionRepository)
	mockLogger := new(mocks.MockLogger)
	receptionService := service.NewReceptionService(mockRepo, mockLogger)

	pvzID := uuid.New()
	receptionID := uuid.New()

	mockRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(receptionID, nil)
	mockRepo.On("CloseReception", mock.Anything, pvzID).Return(nil)
	mockLogger.On("Infow", "Attempting to close reception", "pvzId", pvzID)
	mockLogger.On("Infow", "Reception closed successfully", "pvzId", pvzID)

	// Act
	err := receptionService.CloseReception(context.Background(), pvzID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCloseReception_GetInProgressError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockReceptionRepository)
	mockLogger := new(mocks.MockLogger)
	receptionService := service.NewReceptionService(mockRepo, mockLogger)

	pvzID := uuid.New()
	expectedError := errors.New("database error")

	mockRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(uuid.Nil, expectedError)
	mockLogger.On("Infow", "Attempting to close reception", "pvzId", pvzID)
	mockLogger.On("Errorw", "Failed to get in-progress reception", "pvzId", pvzID, "error", expectedError)

	// Act
	err := receptionService.CloseReception(context.Background(), pvzID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reception lookup failed")
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCloseReception_NoActiveReception(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockReceptionRepository)
	mockLogger := new(mocks.MockLogger)
	receptionService := service.NewReceptionService(mockRepo, mockLogger)

	pvzID := uuid.New()

	mockRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(uuid.Nil, nil)
	mockLogger.On("Infow", "Attempting to close reception", "pvzId", pvzID)
	mockLogger.On("Warnw", "No active reception found", "pvzId", pvzID)

	// Act
	err := receptionService.CloseReception(context.Background(), pvzID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active reception found")
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCloseReception_CloseError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockReceptionRepository)
	mockLogger := new(mocks.MockLogger)
	receptionService := service.NewReceptionService(mockRepo, mockLogger)

	pvzID := uuid.New()
	receptionID := uuid.New()
	expectedError := errors.New("close error")

	mockRepo.On("GetInProgressReception", mock.Anything, pvzID).Return(receptionID, nil)
	mockRepo.On("CloseReception", mock.Anything, pvzID).Return(expectedError)
	mockLogger.On("Infow", "Attempting to close reception", "pvzId", pvzID)
	mockLogger.On("Errorw", "Failed to close reception", "pvzId", pvzID, "error", expectedError)

	// Act
	err := receptionService.CloseReception(context.Background(), pvzID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to close reception")
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
