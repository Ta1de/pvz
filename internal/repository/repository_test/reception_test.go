package repository_test

import (
	"context"
	"database/sql"
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

func TestCreateReception_Success(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewReceptionPostgres(sqlx.NewDb(db, "sqlmock"), mockLogger)

	pvzId := uuid.New()
	expectedReception := model.Reception{
		Id:       uuid.New(),
		DateTime: time.Now().UTC().Truncate(time.Microsecond),
		PvzId:    pvzId,
		Status:   "in_progress",
	}

	mockLogger.On("Infow", "Creating new reception", "pvzId", pvzId).Return()
	mockLogger.On("Infow", "Successfully created reception", "receptionId", expectedReception.Id).Return()

	query := `
		INSERT INTO reception (pvzId, status)
		VALUES ($1, 'in_progress')
		RETURNING id, dateTime, pvzId, status;
	`

	rows := sqlmock.NewRows([]string{"id", "datetime", "pvzid", "status"}).
		AddRow(expectedReception.Id, expectedReception.DateTime, expectedReception.PvzId, expectedReception.Status)

	mockDB.ExpectQuery(query).
		WithArgs(pvzId).
		WillReturnRows(rows)

	result, err := repo.CreateReception(context.Background(), pvzId)

	assert.NoError(t, err)
	assert.Equal(t, expectedReception, result)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestCreateReception_DBError(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewReceptionPostgres(sqlx.NewDb(db, "sqlmock"), mockLogger)

	pvzId := uuid.New()
	dbErr := errors.New("insert failed")

	mockLogger.On("Infow", "Creating new reception", "pvzId", pvzId).Return()
	mockLogger.On("Errorw", "Failed to create reception", "pvzId", pvzId, "error", dbErr).Return()

	query := `
		INSERT INTO reception (pvzId, status)
		VALUES ($1, 'in_progress')
		RETURNING id, dateTime, pvzId, status;
	`

	mockDB.ExpectQuery(query).
		WithArgs(pvzId).
		WillReturnError(dbErr)

	result, err := repo.CreateReception(context.Background(), pvzId)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating reception")
	assert.Equal(t, model.Reception{}, result)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestGetInProgressReception_Success(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewReceptionPostgres(sqlx.NewDb(db, "sqlmock"), mockLogger)

	pvzId := uuid.New()
	expectedReceptionId := uuid.New()

	mockLogger.On("Infow", "Found in-progress reception", "pvzId", pvzId, "receptionId", expectedReceptionId).Return()

	query := `
		SELECT id FROM reception
		WHERE pvzId = $1 AND status = 'in_progress'
		ORDER BY dateTime DESC
		LIMIT 1;
	`

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(expectedReceptionId)

	mockDB.ExpectQuery(query).
		WithArgs(pvzId).
		WillReturnRows(rows)

	result, err := repo.GetInProgressReception(context.Background(), pvzId)

	assert.NoError(t, err)
	assert.Equal(t, expectedReceptionId, result)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestGetInProgressReception_NoRows(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewReceptionPostgres(sqlx.NewDb(db, "sqlmock"), mockLogger)

	pvzId := uuid.New()

	mockLogger.On("Warnw", "No in-progress reception found", "pvzId", pvzId).Return()

	query := `
		SELECT id FROM reception
		WHERE pvzId = $1 AND status = 'in_progress'
		ORDER BY dateTime DESC
		LIMIT 1;
	`

	mockDB.ExpectQuery(query).
		WithArgs(pvzId).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetInProgressReception(context.Background(), pvzId)

	assert.NoError(t, err)
	assert.Equal(t, uuid.Nil, result)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestGetInProgressReception_DBError(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewReceptionPostgres(sqlx.NewDb(db, "sqlmock"), mockLogger)

	pvzId := uuid.New()
	dbErr := errors.New("db is down")

	mockLogger.On("Errorw", "Failed to get in-progress reception", "pvzId", pvzId, "error", dbErr).Return()

	query := `
		SELECT id FROM reception
		WHERE pvzId = $1 AND status = 'in_progress'
		ORDER BY dateTime DESC
		LIMIT 1;
	`

	mockDB.ExpectQuery(query).
		WithArgs(pvzId).
		WillReturnError(dbErr)

	result, err := repo.GetInProgressReception(context.Background(), pvzId)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "get in-progress reception failed")
	assert.Equal(t, uuid.Nil, result)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestCloseReception_Success(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewReceptionPostgres(sqlx.NewDb(db, "sqlmock"), mockLogger)
	pvzId := uuid.New()

	mockLogger.On("Infow", "Successfully closed reception(s)", "pvzId", pvzId).Return()

	query := `
		UPDATE reception
		SET status = 'close'
		WHERE pvzId = $1 AND status = 'in_progress'
	`

	mockDB.ExpectExec(query).
		WithArgs(pvzId).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err = repo.CloseReception(context.Background(), pvzId)

	assert.NoError(t, err)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestCloseReception_DBError(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewReceptionPostgres(sqlx.NewDb(db, "sqlmock"), mockLogger)
	pvzId := uuid.New()
	dbErr := errors.New("connection lost")

	mockLogger.On("Errorw", "Failed to close reception", "pvzId", pvzId, "error", dbErr).Return()

	query := `
		UPDATE reception
		SET status = 'close'
		WHERE pvzId = $1 AND status = 'in_progress'
	`

	mockDB.ExpectExec(query).
		WithArgs(pvzId).
		WillReturnError(dbErr)

	err = repo.CloseReception(context.Background(), pvzId)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to close reception")
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestGetReceptionsByPvzID_Success(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewReceptionPostgres(sqlxDB, mockLogger)

	pvzId := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "datetime", "pvzid", "status"}).
		AddRow(uuid.New(), time.Now(), pvzId, "in_progress").
		AddRow(uuid.New(), time.Now(), pvzId, "close")

	query := `SELECT id, dateTime, pvzId, status FROM reception WHERE pvzId = \$1`

	mockLogger.On("Infow", "Executing GetReceptionsByPvzID query", "pvzId", pvzId).Return()
	mockLogger.On("Infow", "Successfully retrieved receptions list", "count", 2, "pvzId", pvzId).Return()

	mockDB.ExpectQuery(query).
		WithArgs(pvzId).
		WillReturnRows(rows)

	result, err := repo.GetReceptionsByPvzID(context.Background(), pvzId)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestGetReceptionsByPvzID_Error(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewReceptionPostgres(sqlxDB, mockLogger)

	pvzId := uuid.New()
	dbErr := errors.New("query failed")

	query := `SELECT id, dateTime, pvzId, status FROM reception WHERE pvzId = \$1`

	mockLogger.On("Infow", "Executing GetReceptionsByPvzID query", "pvzId", pvzId).Return()
	mockLogger.On("Errorw", "Failed to execute query in GetReceptionsByPvzID", "error", dbErr, "pvzId", pvzId).Return()

	mockDB.ExpectQuery(query).
		WithArgs(pvzId).
		WillReturnError(dbErr)

	result, err := repo.GetReceptionsByPvzID(context.Background(), pvzId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}
