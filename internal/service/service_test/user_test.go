package service_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"pvz/internal/repository/model"
	"pvz/internal/service"
	"pvz/mocks"
)

func TestCreateUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserPostgres)
	mockLogger := new(mocks.MockLogger)

	service := service.NewUserService(mockRepo, mockLogger)
	testUser := model.User{
		Email:    "test@example.com",
		Password: "password123",
	}
	expectedID := uuid.New()

	mockLogger.On("Infow", "User successfully created",
		"userID", expectedID,
		"email", testUser.Email).Once()

	// Ожидаем, что репозиторий вернет ID
	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user model.User) bool {
		// Проверяем, что пароль был захэширован
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
		return err == nil
	})).Return(expectedID, nil)

	// Act
	result, err := service.CreateUser(context.Background(), testUser)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedID, result.Id)
	assert.Equal(t, testUser.Email, result.Email)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateUser_PasswordHashingError(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserPostgres)
	mockLogger := new(mocks.MockLogger)
	service := service.NewUserService(mockRepo, mockLogger)

	// Пароль, который вызовет ошибку хэширования (слишком длинный)
	invalidPassword := string(make([]byte, 100))
	testUser := model.User{
		Email:    "test@example.com",
		Password: invalidPassword,
	}

	// Ожидаем вызов логгера с ошибкой
	mockLogger.On("Errorw", "Password hashing failed",
		"error", mock.Anything).Once()

	// Act
	result, err := service.CreateUser(context.Background(), testUser)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, model.User{}, result)
	assert.Contains(t, err.Error(), "could not hash password")

	mockLogger.AssertExpectations(t)
}

func TestCreateUser_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockUserPostgres)
	mockLogger := new(mocks.MockLogger)
	service := service.NewUserService(mockRepo, mockLogger)

	testUser := model.User{
		Email:    "test@example.com",
		Password: "password123",
	}
	expectedError := errors.New("repository error")

	// Настроим mock репозитория чтобы он возвращал ошибку
	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(uuid.Nil, expectedError)

	// Act
	result, err := service.CreateUser(context.Background(), testUser)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, model.User{}, result)

	mockRepo.AssertExpectations(t)
}

func TestLoginUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserPostgres)
	mockLogger := new(mocks.MockLogger)
	service := service.NewUserService(mockRepo, mockLogger)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	expectedUser := model.User{
		Id:       uuid.New(),
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     "user",
	}

	mockRepo.On("GetUserByEmail", mock.Anything, expectedUser.Email).Return(expectedUser, nil)
	mockLogger.On("Infow", "User authentication successful",
		"userID", expectedUser.Id,
		"email", expectedUser.Email).Once()

	os.Setenv("SIGNING_KEY", "secret-key")

	// Act
	token, err := service.LoginUser(context.Background(), expectedUser.Email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestLoginUser_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserPostgres)
	mockLogger := new(mocks.MockLogger)
	service := service.NewUserService(mockRepo, mockLogger)

	email := "nonexistent@example.com"
	mockRepo.On("GetUserByEmail", mock.Anything, email).Return(model.User{}, errors.New("user not found"))

	// Act
	token, err := service.LoginUser(context.Background(), email, "any-password")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get user")
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
}

func TestLoginUser_WrongPassword(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserPostgres)
	mockLogger := new(mocks.MockLogger)
	service := service.NewUserService(mockRepo, mockLogger)

	correctPassword := "correct-password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

	user := model.User{
		Id:       uuid.New(),
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	mockRepo.On("GetUserByEmail", mock.Anything, user.Email).Return(user, nil)
	mockLogger.On("Warnw", "Incorrect password attempt", "email", user.Email).Once()

	// Act
	token, err := service.LoginUser(context.Background(), user.Email, "wrong-password")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDummyLogin(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserPostgres)
	mockLogger := new(mocks.MockLogger)
	service := service.NewUserService(mockRepo, mockLogger)

	// Set up environment variable for signing key
	originalSigningKey := os.Getenv("SIGNING_KEY")
	defer func() {
		os.Setenv("SIGNING_KEY", originalSigningKey) // Restore original value
	}()
	os.Setenv("SIGNING_KEY", "test-signing-key")

	testRole := "admin"

	// Expected logger call
	mockLogger.On("Infow", "Dummy token created", "role", testRole).Once()

	// Act
	token, err := service.DummyLogin(context.Background(), testRole)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the token can be parsed and contains the correct claims
	parsedToken, err := jwt.ParseWithClaims(token, &model.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-signing-key"), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(*model.TokenClaims)
	assert.True(t, ok)
	assert.Equal(t, testRole, claims.Role)

	mockLogger.AssertExpectations(t)
}

//func TestDummyLogin_SigningError(t *testing.T) {
//	// Arrange
//	mockRepo := new(mocks.MockUserPostgres)
//	mockLogger := new(mocks.MockLogger)
//	service := service.NewUserService(mockRepo, mockLogger)
//
//	// Set empty signing key to force an error
//	originalSigningKey := os.Getenv("SIGNING_KEY")
//	defer func() {
//		os.Setenv("SIGNING_KEY", originalSigningKey) // Restore original value
//	}()
//	os.Setenv("SIGNING_KEY", "")
//
//	testRole := "admin"
//
//	// Expected logger call
//	mockLogger.On("Errorw", "Failed to sign dummy token",
//		"role", testRole,
//		"error", mock.Anything).Once()
//
//	// Act
//	token, err := service.DummyLogin(context.Background(), testRole)
//
//	// Assert
//	assert.Error(t, err)
//	assert.Contains(t, err.Error(), "could not sign token")
//	assert.Empty(t, token)
//
//	mockLogger.AssertExpectations(t)
//}
