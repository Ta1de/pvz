package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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

func TestHandler_CreateReception_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceptionService := mocks.NewMockReception(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Reception: mockReceptionService}
	h := handler.NewHandler(services, mockLogger)

	// Test data
	pvzID := uuid.New()
	reqBody := response.ReceptionRequest{
		PvzId: pvzID.String(),
	}
	jsonBody, _ := json.Marshal(reqBody)

	expectedReception := model.Reception{
		Id:       uuid.New(),
		PvzId:    pvzID,
		DateTime: time.Now(),
		Status:   "open",
	}

	// Mock expectations
	mockReceptionService.EXPECT().
		CreateReception(gomock.Any(), pvzID).
		Return(expectedReception, nil)

	mockLogger.On("Infow", "Reception created successfully",
		"receptionId", expectedReception.Id,
		"pvzId", expectedReception.PvzId).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewBuffer(jsonBody))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.CreateReception(ctx)

	// Verify
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp response.ReceptionResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	respID, err := uuid.Parse(resp.Id)
	assert.NoError(t, err)
	assert.Equal(t, expectedReception.Id, respID)

	respPvzID, err := uuid.Parse(resp.PvzId)
	assert.NoError(t, err)
	assert.Equal(t, expectedReception.PvzId, respPvzID)

	_, err = time.Parse("2006-01-02 15:04:05", resp.DateTime)
	assert.NoError(t, err)
	assert.Equal(t, expectedReception.Status, resp.Status)

	mockLogger.AssertExpectations(t)
}

func TestHandler_CreateReception_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceptionService := mocks.NewMockReception(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Reception: mockReceptionService}
	h := handler.NewHandler(services, mockLogger)

	// Test data - invalid JSON
	invalidJSON := "{invalid}"

	// Mock expectations

	mockLogger.On("Warnw", "Invalid input data for reception creation", "error", mock.Anything).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewBufferString(invalidJSON))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.CreateReception(ctx)

	// Verify
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["message"], "Invalid input data")
	assert.Contains(t, resp, "error")

	mockLogger.AssertExpectations(t)
}

func TestHandler_CreateReception_InvalidPvzId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceptionService := mocks.NewMockReception(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Reception: mockReceptionService}
	h := handler.NewHandler(services, mockLogger)

	// Test data - invalid PvzId format
	reqBody := response.ReceptionRequest{
		PvzId: "invalid-uuid",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Mock expectations
	mockLogger.On("Warnw", "Invalid input data for reception creation", "error", mock.Anything).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewBuffer(jsonBody))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.CreateReception(ctx)

	// Verify
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["message"], "Invalid input data")
	assert.Contains(t, resp, "error")

	mockLogger.AssertExpectations(t)
}

func TestHandler_CreateReception_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceptionService := mocks.NewMockReception(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Reception: mockReceptionService}
	h := handler.NewHandler(services, mockLogger)

	// Test data
	pvzID := uuid.New()
	reqBody := response.ReceptionRequest{
		PvzId: pvzID.String(),
	}
	jsonBody, _ := json.Marshal(reqBody)
	expectedErr := errors.New("database error")

	// Mock expectations
	mockReceptionService.EXPECT().
		CreateReception(gomock.Any(), pvzID).
		Return(model.Reception{}, expectedErr)

	mockLogger.On("Errorw", "Failed to create reception", "error", expectedErr).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewBuffer(jsonBody))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.CreateReception(ctx)

	// Verify
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to create reception", resp["message"])
	assert.Contains(t, resp["error"], expectedErr.Error())

	mockLogger.AssertExpectations(t)
}

func TestHandler_CloseReception_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceptionService := mocks.NewMockReception(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Reception: mockReceptionService}
	h := handler.NewHandler(services, mockLogger)

	// Test data
	pvzID := uuid.New()

	// Mock expectations
	mockReceptionService.EXPECT().
		CloseReception(gomock.Any(), pvzID).
		Return(nil)

	mockLogger.On("Infow", "Attempting to close reception", "PvzId", pvzID).Once()
	mockLogger.On("Infow", "Reception closed successfully", "PvzId", pvzID).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/receptions/"+pvzID.String()+"/close", nil)
	ctx.Params = gin.Params{gin.Param{Key: "pvzId", Value: pvzID.String()}}

	h.CloseReception(ctx)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Reception closed successfully", resp["message"])

	mockLogger.AssertExpectations(t)
}

func TestHandler_CloseReception_InvalidPvzId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceptionService := mocks.NewMockReception(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Reception: mockReceptionService}
	h := handler.NewHandler(services, mockLogger)

	// Test data - invalid PvzId format
	invalidPvzId := "invalid-uuid"

	// Mock expectations
	mockLogger.On("Errorw", "Invalid PvzId format", "PvzId", invalidPvzId, "error", mock.Anything).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/receptions/"+invalidPvzId+"/close", nil)
	ctx.Params = gin.Params{gin.Param{Key: "pvzId", Value: invalidPvzId}}

	h.CloseReception(ctx)

	// Verify
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid PvzId format", resp["error"])

	mockLogger.AssertExpectations(t)
}

func TestHandler_CloseReception_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceptionService := mocks.NewMockReception(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Reception: mockReceptionService}
	h := handler.NewHandler(services, mockLogger)

	// Test data
	pvzID := uuid.New()
	expectedErr := errors.New("reception not found")

	// Mock expectations
	mockReceptionService.EXPECT().
		CloseReception(gomock.Any(), pvzID).
		Return(expectedErr)

	mockLogger.On("Infow", "Attempting to close reception", "PvzId", pvzID).Once()
	mockLogger.On("Errorw", "Failed to close reception", "PvzId", pvzID, "error", expectedErr).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/receptions/"+pvzID.String()+"/close", nil)
	ctx.Params = gin.Params{gin.Param{Key: "pvzId", Value: pvzID.String()}}

	h.CloseReception(ctx)

	// Verify
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], expectedErr.Error())

	mockLogger.AssertExpectations(t)
}
