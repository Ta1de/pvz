package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"pvz/internal/repository"
	"pvz/internal/repository/model"
	"pvz/mocks"
)

func TestCreatePvz_Success(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPvzPostgres(sqlx.NewDb(db, "sqlmock"), mockLogger)

	city := "Moscow"
	expectedPvz := model.Pvz{
		Id:               uuid.New(),
		City:             city,
		RegistrationDate: time.Now().UTC().Truncate(time.Microsecond),
	}

	mockLogger.On("Infow",
		"Inserting new PVZ into database",
		"city", city).Return()

	mockLogger.On("Infow",
		"Successfully inserted PVZ",
		"pvz", expectedPvz).Return()

	exactQuery := `
		INSERT INTO pvz (city)
		VALUES ($1)
		RETURNING id, city, registrationDate
	`

	rows := sqlmock.NewRows([]string{"id", "city", "registrationdate"}).
		AddRow(expectedPvz.Id, expectedPvz.City, expectedPvz.RegistrationDate)

	mockDB.ExpectQuery(exactQuery).
		WithArgs(city).
		WillReturnRows(rows)

	result, err := repo.CreatePvz(context.Background(), city)

	assert.NoError(t, err)
	assert.Equal(t, expectedPvz, result)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestCreatePvz_DBError(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPvzPostgres(sqlx.NewDb(db, "sqlmock"), mockLogger)

	city := "Moscow"
	dbError := errors.New("database error")

	mockLogger.On("Infow",
		"Inserting new PVZ into database",
		"city", city).Return()

	mockLogger.On("Errorw",
		"Failed to insert PVZ",
		"city", city,
		"error", dbError).Return()

	exactQuery := `
		INSERT INTO pvz (city)
		VALUES ($1)
		RETURNING id, city, registrationDate
	`

	mockDB.ExpectQuery(exactQuery).
		WithArgs(city).
		WillReturnError(dbError)

	result, err := repo.CreatePvz(context.Background(), city)

	assert.Error(t, err)
	assert.Equal(t, model.Pvz{}, result)
	assert.Contains(t, err.Error(), "error creating PVZ")
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestGetPvzListByReceptionDate_Success(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewPvzPostgres(sqlxDB, mockLogger)

	// Подготовка данных
	startDate := time.Now().Add(-24 * time.Hour).Truncate(time.Microsecond)
	endDate := time.Now().Truncate(time.Microsecond)
	limit := 10
	offset := 0

	expectedPvz := []model.Pvz{
		{
			Id:               uuid.New(),
			City:             "Moscow",
			RegistrationDate: time.Now().Add(-48 * time.Hour).UTC().Truncate(time.Microsecond),
		},
		{
			Id:               uuid.New(),
			City:             "Kazan",
			RegistrationDate: time.Now().Add(-36 * time.Hour).UTC().Truncate(time.Microsecond),
		},
	}

	// Логгирование
	mockLogger.On("Infow", "Executing GetPvzListByReceptionDate query",
		"startDate", &startDate,
		"endDate", &endDate,
		"limit", limit,
		"offset", offset).Return()

	mockLogger.On("Infow", "Successfully retrieved Pvz list", "count", len(expectedPvz)).Return()

	// SQL-запрос
	query := `
		SELECT DISTINCT p.id, p.registrationDate, p.city
		FROM pvz p
		JOIN reception r ON r.pvzId = p.id
		WHERE ($1::timestamp IS NULL OR r.dateTime >= $1)
		  AND ($2::timestamp IS NULL OR r.dateTime <= $2)
		ORDER BY p.registrationDate DESC
		LIMIT $3 OFFSET $4
	`

	rows := sqlmock.NewRows([]string{"id", "registrationdate", "city"}).
		AddRow(expectedPvz[0].Id, expectedPvz[0].RegistrationDate, expectedPvz[0].City).
		AddRow(expectedPvz[1].Id, expectedPvz[1].RegistrationDate, expectedPvz[1].City)

	mockDB.ExpectQuery(query).
		WithArgs(startDate, endDate, limit, offset).
		WillReturnRows(rows)

	// Вызов метода
	result, err := repo.GetPvzListByReceptionDate(context.Background(), limit, offset, &startDate, &endDate)

	assert.NoError(t, err)
	assert.Equal(t, expectedPvz, result)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestGetPvzListByReceptionDate_DBError(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewPvzPostgres(sqlxDB, mockLogger)

	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()
	limit := 5
	offset := 0
	dbError := errors.New("db error")

	// Логгирование
	mockLogger.On("Infow", "Executing GetPvzListByReceptionDate query",
		"startDate", &startDate,
		"endDate", &endDate,
		"limit", limit,
		"offset", offset).Return()

	mockLogger.On("Errorw", "Failed to fetch Pvz list", "error", dbError).Return()

	// SQL-запрос
	query := `
		SELECT DISTINCT p.id, p.registrationDate, p.city
		FROM pvz p
		JOIN reception r ON r.pvzId = p.id
		WHERE ($1::timestamp IS NULL OR r.dateTime >= $1)
		  AND ($2::timestamp IS NULL OR r.dateTime <= $2)
		ORDER BY p.registrationDate DESC
		LIMIT $3 OFFSET $4
	`

	mockDB.ExpectQuery(query).
		WithArgs(startDate, endDate, limit, offset).
		WillReturnError(dbError)

	result, err := repo.GetPvzListByReceptionDate(context.Background(), limit, offset, &startDate, &endDate)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "db error")
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}
