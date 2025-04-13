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
	"github.com/stretchr/testify/mock"
	"pvz/internal/repository"
	"pvz/internal/repository/model"
	"pvz/mocks"
)

func TestCreateProduct_Success(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(sqlx.NewDb(db, "sqlmock"), mockLogger)

	product := model.Product{
		Type:        "TestType",
		ReceptionId: uuid.New(),
	}

	mockLogger.On("Infow",
		"Successfully created product",
		"product", mock.AnythingOfType("model.Product")).Return()

	rows := sqlmock.NewRows([]string{"id", "datetime", "type", "receptionid"}).
		AddRow(uuid.New(), time.Now(), product.Type, product.ReceptionId)

	mockDB.ExpectQuery(`INSERT INTO product`).
		WithArgs(product.Type, product.ReceptionId).
		WillReturnRows(rows)

	result, err := repo.CreateProduct(context.Background(), product)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, result.Id)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestCreateProduct_Error(t *testing.T) {

	mockLogger := new(mocks.MockLogger)

	db, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewRepository(sqlxDB, mockLogger)

	testUUID := uuid.New()
	productType := "TestType"

	exactQuery := `
	INSERT INTO product (type, receptionid)
	VALUES ($1, $2)
	RETURNING id, datetime, type, receptionid;
	`

	mockLogger.On("Errorw",
		"Failed to create product",
		"product", mock.Anything,
		"error", mock.Anything).Return()

	mockDB.ExpectQuery(exactQuery).
		WithArgs(productType, testUUID).
		WillReturnError(errors.New("database error"))

	_, err = repo.CreateProduct(context.Background(), model.Product{
		Type:        productType,
		ReceptionId: testUUID,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error inserting product")
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestGetLastProductIdByReception_Success(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(sqlx.NewDb(db, "sqlmock"), mockLogger)

	receptionId := uuid.New()
	expectedId := uuid.New()

	mockLogger.On("Infow",
		"Fetched last product ID",
		"productId", expectedId,
		"receptionId", receptionId).Return()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedId)

	mockDB.ExpectQuery(`SELECT id FROM product WHERE receptionId = \$1 ORDER BY datetime DESC LIMIT 1`).
		WithArgs(receptionId).
		WillReturnRows(rows)

	result, err := repo.GetLastProductIdByReception(context.Background(), receptionId)

	assert.NoError(t, err)
	assert.Equal(t, expectedId, result)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestGetLastProductIdByReception_NoRows(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(sqlx.NewDb(db, "sqlmock"), mockLogger)

	receptionId := uuid.New()

	mockLogger.On("Warnw",
		"No products found for reception",
		"receptionId", receptionId).Return()

	mockDB.ExpectQuery(`SELECT id FROM product WHERE receptionId = \$1 ORDER BY datetime DESC LIMIT 1`).
		WithArgs(receptionId).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetLastProductIdByReception(context.Background(), receptionId)

	assert.NoError(t, err)
	assert.Equal(t, uuid.Nil, result)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestGetLastProductIdByReception_DBError(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(sqlx.NewDb(db, "sqlmock"), mockLogger)

	receptionId := uuid.New()
	dbError := errors.New("database error")

	mockLogger.On("Errorw",
		"Failed to fetch last product ID",
		"receptionId", receptionId,
		"error", dbError).Return()

	mockDB.ExpectQuery(`SELECT id FROM product WHERE receptionId = \$1 ORDER BY datetime DESC LIMIT 1`).
		WithArgs(receptionId).
		WillReturnError(dbError)

	result, err := repo.GetLastProductIdByReception(context.Background(), receptionId)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, result)
	assert.Contains(t, err.Error(), "failed to get last product id")
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestDeleteProductById_Success(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(sqlx.NewDb(db, "sqlmock"), mockLogger)

	productId := uuid.New()

	mockLogger.On("Infow",
		"Product deleted successfully",
		"productId", productId).Return()

	mockDB.ExpectExec(`DELETE FROM product WHERE id = \$1`).
		WithArgs(productId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteProductById(context.Background(), productId)

	assert.NoError(t, err)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestDeleteProductById_DBError(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(sqlx.NewDb(db, "sqlmock"), mockLogger)

	productId := uuid.New()
	dbError := errors.New("database error")

	mockLogger.On("Errorw",
		"Failed to delete product",
		"productId", productId,
		"error", dbError).Return()

	mockDB.ExpectExec(`DELETE FROM product WHERE id = \$1`).
		WithArgs(productId).
		WillReturnError(dbError)

	err = repo.DeleteProductById(context.Background(), productId)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete product")
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestDeleteProductById_NoRowsAffected(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(sqlx.NewDb(db, "sqlmock"), mockLogger)

	productId := uuid.New()

	mockLogger.On("Infow",
		"Product deleted successfully",
		"productId", productId).Return()

	mockDB.ExpectExec(`DELETE FROM product WHERE id = \$1`).
		WithArgs(productId).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.DeleteProductById(context.Background(), productId)

	assert.NoError(t, err)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestGetProductsByReceptionID_Success(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(sqlx.NewDb(db, "sqlmock"), mockLogger)

	receptionId := uuid.New()
	expectedProducts := []model.Product{
		{
			Id:          uuid.New(),
			DateTime:    time.Now().UTC().Truncate(time.Microsecond),
			Type:        "Type1",
			ReceptionId: receptionId,
		},
		{
			Id:          uuid.New(),
			DateTime:    time.Now().UTC().Truncate(time.Microsecond).Add(time.Hour),
			Type:        "Type2",
			ReceptionId: receptionId,
		},
	}

	mockLogger.On("Infow",
		"Executing GetProductsByReceptionID query",
		"receptionId", receptionId).Return()

	mockLogger.On("Infow",
		"Successfully retrieved products",
		"count", len(expectedProducts),
		"receptionId", receptionId).Return()

	rows := sqlmock.NewRows([]string{"id", "datetime", "type", "receptionid"})
	for _, p := range expectedProducts {
		rows.AddRow(p.Id, p.DateTime, p.Type, p.ReceptionId)
	}

	mockDB.ExpectQuery(`SELECT id, datetime, type, receptionId FROM product WHERE receptionId = \$1`).
		WithArgs(receptionId).
		WillReturnRows(rows)

	result, err := repo.GetProductsByReceptionID(context.Background(), receptionId)

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, result)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestGetProductsByReceptionID_EmptyResult(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(sqlx.NewDb(db, "sqlmock"), mockLogger)

	receptionId := uuid.New()

	mockLogger.On("Infow",
		"Executing GetProductsByReceptionID query",
		"receptionId", receptionId).Return()

	mockLogger.On("Infow",
		"Successfully retrieved products",
		"count", 0,
		"receptionId", receptionId).Return()

	rows := sqlmock.NewRows([]string{"id", "datetime", "type", "receptionid"})

	mockDB.ExpectQuery(`SELECT id, datetime, type, receptionId FROM product WHERE receptionId = \$1`).
		WithArgs(receptionId).
		WillReturnRows(rows)

	result, err := repo.GetProductsByReceptionID(context.Background(), receptionId)

	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}

func TestGetProductsByReceptionID_DBError(t *testing.T) {
	mockLogger := new(mocks.MockLogger)
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewRepository(sqlx.NewDb(db, "sqlmock"), mockLogger)

	receptionId := uuid.New()
	dbError := errors.New("database error")

	mockLogger.On("Infow",
		"Executing GetProductsByReceptionID query",
		"receptionId", receptionId).Return()

	mockLogger.On("Errorw",
		"Failed to fetch products",
		"error", dbError,
		"receptionId", receptionId).Return()

	mockDB.ExpectQuery(`SELECT id, datetime, type, receptionId FROM product WHERE receptionId = \$1`).
		WithArgs(receptionId).
		WillReturnError(dbError)

	result, err := repo.GetProductsByReceptionID(context.Background(), receptionId)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, dbError, err)
	assert.NoError(t, mockDB.ExpectationsWereMet())
	mockLogger.AssertExpectations(t)
}
