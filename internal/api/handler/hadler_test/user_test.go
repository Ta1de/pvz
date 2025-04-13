package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"pvz/internal/api/handler"
	"pvz/internal/api/response"
	"pvz/internal/repository/model"
	"pvz/internal/service"
	"pvz/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/mock/gomock"
)

func TestHandler_DummyLogin_Success(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUser(ctrl)
	mockLogger := new(mocks.MockLogger)

	h := handler.NewHandler(&service.Service{User: mockService}, mockLogger)

	// Test data
	role := "admin"
	token := "test-token"
	reqBody := response.DummyLoginPostRequest{Role: role}
	jsonBody, _ := json.Marshal(reqBody)

	// Mock expectations
	mockService.EXPECT().
		DummyLogin(gomock.Any(), role).
		Return(token, nil)

	mockLogger.On("Infow", "Dummy login successful", "role", role).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/dummy-login", bytes.NewBuffer(jsonBody))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.DummyLogin(ctx)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, token, resp["token"])

	mockLogger.AssertExpectations(t)
}

func TestHandler_DummyLogin_InvalidInput(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUser(ctrl)
	mockLogger := new(mocks.MockLogger)

	h := handler.NewHandler(&service.Service{User: mockService}, mockLogger)

	// Test data - invalid JSON
	invalidJSON := "{invalid}"

	// Mock expectations
	mockLogger.On("Warnw", "Invalid input data for dummy login", "error", mock.Anything).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/dummy-login", bytes.NewBufferString(invalidJSON))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.DummyLogin(ctx)

	// Verify
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["message"], "Invalid input data")
	assert.Contains(t, resp, "error")

	mockLogger.AssertExpectations(t)
}

func TestHandler_DummyLogin_ServiceError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUser(ctrl)
	mockLogger := new(mocks.MockLogger)

	h := handler.NewHandler(&service.Service{User: mockService}, mockLogger)

	// Test data
	role := "admin"
	reqBody := response.DummyLoginPostRequest{Role: role}
	jsonBody, _ := json.Marshal(reqBody)
	expectedErr := errors.New("service error")

	// Mock expectations
	mockService.EXPECT().
		DummyLogin(gomock.Any(), role).
		Return("", expectedErr)

	mockLogger.On("Warnw", "Dummy login failed", "error", expectedErr).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/dummy-login", bytes.NewBuffer(jsonBody))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.DummyLogin(ctx)

	// Verify
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to generate token", resp["message"])

	mockLogger.AssertExpectations(t)
}

func TestHandler_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUser(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{User: mockUserService}
	h := handler.NewHandler(services, mockLogger)

	// Test data
	reqBody := response.LoginPostRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)
	token := "test-token"

	// Mock expectations
	mockUserService.EXPECT().
		LoginUser(gomock.Any(), reqBody.Email, reqBody.Password).
		Return(token, nil)

	mockLogger.On("Infow", "Login successful", "email", reqBody.Email).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.Login(ctx)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, token, resp["token"])

	mockLogger.AssertExpectations(t)
}

func TestHandler_Login_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUser(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{User: mockUserService}
	h := handler.NewHandler(services, mockLogger)

	// Test data - invalid JSON
	invalidJSON := "{invalid}"

	// Mock expectations
	mockLogger.On("Warnw", "Invalid input data for login", "error", mock.Anything).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(invalidJSON))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.Login(ctx)

	// Verify
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["message"], "Invalid input data")
	assert.Contains(t, resp, "error")

	mockLogger.AssertExpectations(t)
}

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUser(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{User: mockUserService}
	h := handler.NewHandler(services, mockLogger)

	// Test data
	reqBody := response.LoginPostRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	jsonBody, _ := json.Marshal(reqBody)
	expectedErr := errors.New("invalid credentials")

	// Mock expectations
	mockUserService.EXPECT().
		LoginUser(gomock.Any(), reqBody.Email, reqBody.Password).
		Return("", expectedErr)

	mockLogger.On("Warnw", "Login failed", "email", reqBody.Email, "error", expectedErr).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.Login(ctx)

	// Verify
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid email or password", resp["message"])

	mockLogger.AssertExpectations(t)
}

func TestHandler_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUser(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{User: mockUserService}
	h := handler.NewHandler(services, mockLogger)

	// Test data
	reqBody := response.RegisterPostRequest{
		Email:    "test@example.com",
		Password: "password123",
		Role:     "user", // Добавляем роль, так как она есть в структуре запроса
	}
	jsonBody, _ := json.Marshal(reqBody)

	expectedUser := model.User{
		Id:       uuid.New(),
		Email:    reqBody.Email,
		Password: "hashed_password",
		Role:     reqBody.Role,
	}

	// Mock expectations
	mockUserService.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, user model.User) (model.User, error) {
			assert.Equal(t, reqBody.Email, user.Email)
			assert.Equal(t, reqBody.Role, user.Role)
			return expectedUser, nil
		})

	mockLogger.On("Infow", "User registered successfully",
		"userID", expectedUser.Id,
		"email", expectedUser.Email).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.Register(ctx)

	// Verify
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp response.RegisterResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Преобразуем строковый ID обратно в UUID для сравнения
	respUUID, err := uuid.Parse(resp.Id)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Id, respUUID)
	assert.Equal(t, expectedUser.Email, resp.Email)
	assert.Equal(t, expectedUser.Role, resp.Role)

	mockLogger.AssertExpectations(t)
}

func TestHandler_Register_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUser(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{User: mockUserService}
	h := handler.NewHandler(services, mockLogger)

	// Test data - invalid JSON
	invalidJSON := "{invalid}"

	// Mock expectations
	mockLogger.On("Warnw", "Invalid input data for registration", "error", mock.Anything).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(invalidJSON))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.Register(ctx)

	// Verify
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["message"], "Invalid input data")
	assert.Contains(t, resp, "error")

	mockLogger.AssertExpectations(t)
}

func TestHandler_Register_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUser(ctrl)
	mockLogger := new(mocks.MockLogger)

	services := &service.Service{User: mockUserService}
	h := handler.NewHandler(services, mockLogger)

	// Test data
	reqBody := response.RegisterPostRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)
	expectedErr := errors.New("database error")

	// Mock expectations
	mockUserService.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(model.User{}, expectedErr)

	mockLogger.On("Errorw", "User registration failed", "error", expectedErr).Once()

	// Execute
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.Register(ctx)

	// Verify
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to create user", resp["message"])
	assert.Contains(t, resp["error"], expectedErr.Error())

	mockLogger.AssertExpectations(t)
}
