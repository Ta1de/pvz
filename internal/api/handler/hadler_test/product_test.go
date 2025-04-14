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

func TestHandler_AddProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductService := mocks.NewMockProduct(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Product: mockProductService}
	h := handler.NewHandler(services, mockLogger)

	// Test data
	pvzID := uuid.New()
	productType := "package"
	reqBody := response.ProductRequest{
		PvzId: pvzID.String(),
		Type:  productType,
	}
	jsonBody, _ := json.Marshal(reqBody)

	expectedProduct := model.Product{
		Id:          uuid.New(),
		Type:        productType,
		DateTime:    time.Now(),
		ReceptionId: uuid.New(),
	}

	// Mock expectations
	mockProductService.EXPECT().
		AddProduct(gomock.Any(), pvzID, productType).
		Return(expectedProduct, nil)

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(jsonBody))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.AddProduct(ctx)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.ProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	respID, err := uuid.Parse(resp.Id)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct.Id, respID)

	_, err = time.Parse("2006-01-02 15:04:05", resp.DateTime)
	assert.NoError(t, err)

	assert.Equal(t, expectedProduct.Type, resp.Type)

	receptionID, err := uuid.Parse(resp.ReceptionId)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct.ReceptionId, receptionID)

	mockLogger.AssertExpectations(t)
}

func TestHandler_AddProduct_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductService := mocks.NewMockProduct(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Product: mockProductService}
	h := handler.NewHandler(services, mockLogger)

	// Test data - invalid JSON
	invalidJSON := "{invalid}"

	// Mock expectations
	mockLogger.On("Errorw", "Failed to bind product request", "error", mock.Anything).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(invalidJSON))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.AddProduct(ctx)

	// Verify
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request body", resp["error"])
	assert.Contains(t, resp, "details")

	mockLogger.AssertExpectations(t)
}

func TestHandler_AddProduct_InvalidPvzId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductService := mocks.NewMockProduct(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Product: mockProductService}
	h := handler.NewHandler(services, mockLogger)

	// Test data - invalid PvzId format
	reqBody := response.ProductRequest{
		PvzId: "invalid-uuid",
		Type:  "package",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Mock expectations
	mockLogger.On("Errorw", "Invalid PvzId format", "PvzId", "invalid-uuid", "error", mock.Anything).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(jsonBody))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.AddProduct(ctx)

	// Verify
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid PvzId format", resp["error"])

	mockLogger.AssertExpectations(t)
}

func TestHandler_AddProduct_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductService := mocks.NewMockProduct(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Product: mockProductService}
	h := handler.NewHandler(services, mockLogger)

	// Test data
	pvzID := uuid.New()
	productType := "package"
	reqBody := response.ProductRequest{
		PvzId: pvzID.String(),
		Type:  productType,
	}
	jsonBody, _ := json.Marshal(reqBody)
	expectedErr := errors.New("storage error")

	// Mock expectations
	mockProductService.EXPECT().
		AddProduct(gomock.Any(), pvzID, productType).
		Return(model.Product{}, expectedErr)

	mockLogger.On("Errorw", "Failed to add product", "error", expectedErr, "PvzId", pvzID, "type", productType).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(jsonBody))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.AddProduct(ctx)

	// Verify
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to add product", resp["error"])

	mockLogger.AssertExpectations(t)
}

func TestHandler_DeleteLastProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductService := mocks.NewMockProduct(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Product: mockProductService}
	h := handler.NewHandler(services, mockLogger)

	// Test data
	pvzID := uuid.New()

	// Mock expectations
	mockProductService.EXPECT().
		DeleteLastProduct(gomock.Any(), pvzID).
		Return(nil)

	mockLogger.On("Infow", "Attempting to delete last product", "PvzId", pvzID).Once()
	mockLogger.On("Infow", "Last product deleted successfully", "PvzId", pvzID).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/products/last/"+pvzID.String(), nil)
	ctx.Params = gin.Params{gin.Param{Key: "pvzId", Value: pvzID.String()}}

	h.DeleteLastProduct(ctx)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Last product deleted successfully", resp["message"])

	mockLogger.AssertExpectations(t)
}

func TestHandler_DeleteLastProduct_InvalidPvzId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductService := mocks.NewMockProduct(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Product: mockProductService}
	h := handler.NewHandler(services, mockLogger)

	// Test data - invalid PvzId format
	invalidPvzId := "invalid-uuid"

	// Mock expectations
	mockLogger.On("Errorw", "Invalid PvzId format", "PvzId", invalidPvzId, "error", mock.Anything).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/products/last/"+invalidPvzId, nil)
	ctx.Params = gin.Params{gin.Param{Key: "pvzId", Value: invalidPvzId}}

	h.DeleteLastProduct(ctx)

	// Verify
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid PvzId format", resp["error"])

	mockLogger.AssertExpectations(t)
}

func TestHandler_DeleteLastProduct_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductService := mocks.NewMockProduct(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{Product: mockProductService}
	h := handler.NewHandler(services, mockLogger)

	// Test data
	pvzID := uuid.New()
	expectedErr := errors.New("no products found")

	// Mock expectations
	mockProductService.EXPECT().
		DeleteLastProduct(gomock.Any(), pvzID).
		Return(expectedErr)

	mockLogger.On("Infow", "Attempting to delete last product", "PvzId", pvzID).Once()
	mockLogger.On("Errorw", "Failed to delete last product", "PvzId", pvzID, "error", expectedErr).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/products/last/"+pvzID.String(), nil)
	ctx.Params = gin.Params{gin.Param{Key: "pvzId", Value: pvzID.String()}}

	h.DeleteLastProduct(ctx)

	// Verify
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], expectedErr.Error())

	mockLogger.AssertExpectations(t)
}
