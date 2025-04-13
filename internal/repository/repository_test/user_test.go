package repository_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pvz/internal/repository"
	"pvz/internal/repository/model"
	"pvz/mocks"
)

func TestCreateUser_Success(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewUserPostgres(sqlxDB, mockLogger)

	user := model.User{
		Email:    "test@example.com",
		Role:     "admin",
		Password: "hashed-password",
	}
	expectedID := uuid.New()

	mock.ExpectQuery(`INSERT INTO users \(email, role, password\)\s+VALUES \(\$1, \$2, \$3\)\s+RETURNING id;`).
		WithArgs(user.Email, user.Role, user.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	mockLogger.On("Infow", "User created in database", "userID", expectedID, "email", user.Email).Return()

	id, err := repo.CreateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
	assert.NoError(t, mock.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestCreateUser_DBError(t *testing.T) {
	// Инициализация моков
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создание репозитория
	repo := repository.NewUserPostgres(sqlx.NewDb(db, "sqlmock"), mockLogger)

	// Тестовые данные
	user := model.User{
		Email:    "fail@example.com",
		Role:     "user",
		Password: "failpass",
	}
	expectedErr := errors.New("database error")

	// Настройка ожиданий
	mockLogger.On("Errorw",
		"Failed to insert user into database",
		"email", user.Email,
		"error", expectedErr).Return()

	mockDB.ExpectQuery(`INSERT INTO users \(email, role, password\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs(user.Email, user.Role, user.Password).
		WillReturnError(expectedErr)

	// Вызов метода
	id, err := repo.CreateUser(context.Background(), user)

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
	assert.Contains(t, err.Error(), "failed to create user")
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestCreateUser_ScanError(t *testing.T) {
	// Инициализация моков
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создание репозитория
	repo := repository.NewUserPostgres(sqlx.NewDb(db, "sqlmock"), mockLogger)

	// Тестовые данные
	user := model.User{
		Email:    "invalid@example.com",
		Role:     "user",
		Password: "pass",
	}

	// Настройка ожиданий
	mockLogger.On("Errorw",
		"Failed to insert user into database",
		"email", user.Email,
		"error", mock.Anything).Return()

	// Возвращаем невалидные данные для сканирования
	mockDB.ExpectQuery(`INSERT INTO users \(email, role, password\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs(user.Email, user.Role, user.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("not-a-uuid"))

	// Вызов метода
	id, err := repo.CreateUser(context.Background(), user)

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
	assert.Contains(t, err.Error(), "failed to create user")
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestGetUserByEmail_Success(t *testing.T) {
	// Инициализация моков
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создание репозитория
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewUserPostgres(sqlxDB, mockLogger)

	// Тестовые данные
	email := "user@example.com"
	expectedUser := model.User{
		Id:       uuid.New(),
		Email:    email,
		Role:     "user",
		Password: "hashed-password",
	}

	// Ожидания для SQL-запроса
	query := `SELECT id, email, role, password FROM users WHERE email = \$1`
	mockDB.ExpectQuery(query).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "role", "password"}).
			AddRow(expectedUser.Id, expectedUser.Email, expectedUser.Role, expectedUser.Password))

	// Вызов метода
	user, err := repo.GetUserByEmail(context.Background(), email)

	// Проверки
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestGetUserByEmail_UserNotFound(t *testing.T) {
	// Инициализация моков
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создание репозитория
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewUserPostgres(sqlxDB, mockLogger)

	// Тестовые данные
	email := "nonexistent@example.com"
	expectedErr := errors.New("user not found")

	// Ожидания для SQL-запроса (ошибка, пользователь не найден)
	query := `SELECT id, email, role, password FROM users WHERE email = \$1`
	mockDB.ExpectQuery(query).
		WithArgs(email).
		WillReturnError(expectedErr)

	// Ожидание логгера
	mockLogger.On("Warnw", "User not found", "email", email, "error", expectedErr).Return()

	// Вызов метода
	user, err := repo.GetUserByEmail(context.Background(), email)

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("user not found: %w", expectedErr), err)
	assert.Equal(t, model.User{}, user) // Пустой пользователь в случае ошибки
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}
