package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"pvz/internal/api/handler"
	"pvz/internal/api/response"
	"pvz/internal/repository/model"
	"pvz/internal/service"
	"pvz/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/mock/gomock"
)

func TestHandler_CreatePvz_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPvzService := mocks.NewMockPvz(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Pvz: mockPvzService}
	h := handler.NewHandler(services, mockLogger)

	// Test data
	reqBody := response.PvzRequest{
		City: "Moscow",
	}
	jsonBody, _ := json.Marshal(reqBody)

	now := time.Now()
	expectedPvz := model.Pvz{
		Id:               uuid.New(),
		RegistrationDate: now,
		City:             reqBody.City,
	}

	// Mock expectations
	mockPvzService.EXPECT().
		CreatePvz(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, pvz model.Pvz) (model.Pvz, error) {
			assert.Equal(t, reqBody.City, pvz.City)
			return expectedPvz, nil
		})

	mockLogger.On("Infow", "Creating new PVZ", "request", reqBody).Once()
	mockLogger.On("Infow", "Successfully created PVZ", "pvz", mock.Anything).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBuffer(jsonBody))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.CreatePvz(ctx)

	// Verify
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp response.PvzResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	respID, err := uuid.Parse(resp.Id)
	assert.NoError(t, err)
	assert.Equal(t, expectedPvz.Id, respID)
	assert.Equal(t, expectedPvz.City, resp.City)

	_, err = time.Parse("2006-01-02 15:04:05", resp.RegistrationDate)
	assert.NoError(t, err)

	mockLogger.AssertExpectations(t)
}

func TestHandler_CreatePvz_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPvzService := mocks.NewMockPvz(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Pvz: mockPvzService}
	h := handler.NewHandler(services, mockLogger)

	// Test data - invalid JSON
	invalidJSON := "{invalid}"

	// Mock expectations
	mockLogger.On("Warnw", "Invalid PvzRequest", "error", mock.Anything).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBufferString(invalidJSON))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.CreatePvz(ctx)

	// Verify
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "неверный формат данных")

	mockLogger.AssertExpectations(t)
}

func TestHandler_CreatePvz_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPvzService := mocks.NewMockPvz(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Pvz: mockPvzService}
	h := handler.NewHandler(services, mockLogger)

	// Test data
	reqBody := response.PvzRequest{
		City: "Moscow",
	}
	jsonBody, _ := json.Marshal(reqBody)
	expectedErr := errors.New("database error")

	// Mock expectations
	mockPvzService.EXPECT().
		CreatePvz(gomock.Any(), gomock.Any()).
		Return(model.Pvz{}, expectedErr)

	mockLogger.On("Infow", "Creating new PVZ", "request", reqBody).Once()
	mockLogger.On("Errorw", "Failed to create PVZ", "error", expectedErr).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBuffer(jsonBody))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.CreatePvz(ctx)

	// Verify
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), expectedErr.Error())

	mockLogger.AssertExpectations(t)
}

func TestHandler_GetPvz_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPvzService := mocks.NewMockPvz(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Pvz: mockPvzService}
	h := handler.NewHandler(services, mockLogger)

	// Test data
	limit := 10
	offset := 0
	startDateStr := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	endDateStr := time.Now().Format(time.RFC3339)

	startDate, _ := time.Parse(time.RFC3339, startDateStr)
	endDate, _ := time.Parse(time.RFC3339, endDateStr)

	expectedResult := []response.PvzFullResponse{
		{
			Pvz: response.PvzResponse{
				Id:               uuid.New().String(),
				RegistrationDate: time.Now().Format(time.RFC3339),
				City:             "Moscow",
			},
			Receptions: []response.ReceptionWrapper{
				{
					Reception: response.ReceptionResponse{
						Id:     uuid.New().String(),
						PvzId:  uuid.New().String(),
						Status: "open",
					},
					Products: []response.ProductResponse{
						{
							Id:          uuid.New().String(),
							ReceptionId: uuid.New().String(),
							Type:        "package",
						},
					},
				},
			},
		},
	}

	// Mock expectations - важно передать указатели на time.Time
	mockPvzService.EXPECT().
		GetPvzList(gomock.Any(), limit, offset, &startDate, &endDate).
		Return(expectedResult, nil)

	mockLogger.On("Infow", "Received request for Pvz list",
		"limit", "10", "offset", "0", "startDate", startDateStr, "endDate", endDateStr).Once()
	mockLogger.On("Infow", "Successfully retrieved Pvz list", "count", len(expectedResult)).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/pvz?limit=10&offset=0&startDate="+url.QueryEscape(startDateStr)+"&endDate="+url.QueryEscape(endDateStr), nil)

	h.GetPvz(ctx)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)

	var resp []response.PvzFullResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, resp)

	mockLogger.AssertExpectations(t)
}

func TestHandler_GetPvz_InvalidLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPvzService := mocks.NewMockPvz(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Pvz: mockPvzService}
	h := handler.NewHandler(services, mockLogger)

	// Mock expectations
	mockLogger.On("Infow", "Received request for Pvz list",
		"limit", "invalid", "offset", "0", "startDate", "", "endDate", "").Once()
	mockLogger.On("Warnw", "Invalid limit", "error", mock.Anything).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/pvz?limit=invalid&offset=0", nil)

	h.GetPvz(ctx)

	// Verify
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid limit")

	mockLogger.AssertExpectations(t)
}

func TestHandler_GetPvz_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPvzService := mocks.NewMockPvz(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Pvz: mockPvzService}
	h := handler.NewHandler(services, mockLogger)

	limit := 10
	offset := 0
	expectedErr := errors.New("database error")

	// Mock expectations
	mockPvzService.EXPECT().
		GetPvzList(gomock.Any(), limit, offset, nil, nil).
		Return(nil, expectedErr)

	mockLogger.On("Infow", "Received request for Pvz list",
		"limit", "10", "offset", "0", "startDate", "", "endDate", "").Once()
	mockLogger.On("Errorw", "Failed to get Pvz list", "error", expectedErr).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/pvz?limit=10&offset=0", nil)

	h.GetPvz(ctx)

	// Verify
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), expectedErr.Error())

	mockLogger.AssertExpectations(t)
}
