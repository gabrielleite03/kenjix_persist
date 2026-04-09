package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/stretchr/testify/assert"
)

func newMock() (*sql.DB, sqlmock.Sqlmock, *StockDAO) {
	db, mock, _ := sqlmock.New()
	dao := &StockDAO{db: db}
	return db, mock, dao
}

func TestStockDAO_Upsert(t *testing.T) {

	db, mock, dao := newMock()
	defer db.Close()

	stock := &model.Stock{
		ProductID:        1,
		WarehousePlaceID: 2,
		Quantity:         10,
		Active:           true,
	}

	query := regexp.QuoteMeta(`
	INSERT INTO stock (product_id, warehouse_place_id, quantity, active)
	VALUES (?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	    quantity = quantity + VALUES(quantity),
	    active = VALUES(active),
	    updated_at = NOW()
	`)

	mock.ExpectExec(query).
		WithArgs(1, 2, 10, true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := dao.Upsert(stock)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStockDAO_Get(t *testing.T) {

	db, mock, dao := newMock()
	defer db.Close()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"product_id",
		"warehouse_place_id",
		"quantity",
		"active",
		"updated_at",
	}).AddRow(1, 2, 10, true, now)

	query := regexp.QuoteMeta(`
	SELECT product_id, warehouse_place_id, quantity, active, updated_at
	FROM stock
	WHERE product_id = ? AND warehouse_place_id = ?
	`)

	mock.ExpectQuery(query).
		WithArgs(1, 2).
		WillReturnRows(rows)

	result, err := dao.Get(1, 2)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ProductID)
	assert.Equal(t, 10, result.Quantity)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStockDAO_InsertMovement(t *testing.T) {

	db, mock, dao := newMock()
	defer db.Close()

	m := &model.StockMovement{
		ProductID:        1,
		WarehousePlaceID: 2,
		Type:             model.StockMovementIn,
		Quantity:         5,
		Reason:           "purchase",
	}

	query := regexp.QuoteMeta(`
	INSERT INTO stock_movement
	(product_id, warehouse_place_id, type, quantity, reference_id, reference_type, reason)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`)

	mock.ExpectExec(query).
		WithArgs(1, 2, model.StockMovementIn, 5, nil, nil, "purchase").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := dao.InsertMovement(m)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStockDAO_AddStock(t *testing.T) {

	db, mock, dao := newMock()
	defer db.Close()

	m := &model.StockMovement{
		ProductID:        1,
		WarehousePlaceID: 2,
		Type:             model.StockMovementIn,
		Quantity:         10,
		Reason:           "purchase",
	}

	insertMovementQuery := regexp.QuoteMeta(`
	INSERT INTO stock_movement
	(product_id, warehouse_place_id, type, quantity, reference_id, reference_type, reason)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`)

	upsertQuery := regexp.QuoteMeta(`
	INSERT INTO stock (product_id, warehouse_place_id, quantity, active)
	VALUES (?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	    quantity = quantity + VALUES(quantity),
	    active = VALUES(active)
	`)

	mock.ExpectBegin()

	mock.ExpectExec(insertMovementQuery).
		WithArgs(1, 2, model.StockMovementIn, 10, nil, nil, "purchase").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(upsertQuery).
		WithArgs(1, 2, 10, true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := dao.AddStock(m)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStockDAO_GetAll(t *testing.T) {

	db, mock, dao := newMock()
	defer db.Close()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"product_id",
		"warehouse_place_id",
		"quantity",
		"active",
		"updated_at",
	}).
		AddRow(1, 1, 10, true, now).
		AddRow(1, 2, 5, true, now)

	query := regexp.QuoteMeta(`
	SELECT product_id, warehouse_place_id, quantity, active, updated_at
	FROM stock
	ORDER BY product_id, warehouse_place_id
	`)

	mock.ExpectQuery(query).
		WillReturnRows(rows)

	result, err := dao.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(1), result[0].ProductID)
	assert.Equal(t, 10, result[0].Quantity)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStockDAO_GetMovementsByProduct(t *testing.T) {

	db, mock, dao := newMock()
	defer db.Close()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"product_id",
		"warehouse_place_id",
		"type",
		"quantity",
		"reference_id",
		"reference_type",
		"reason",
		"created_at",
	}).AddRow(
		1, 2, "IN", 10, nil, nil, "purchase", now,
	)

	query := regexp.QuoteMeta(`
	SELECT 
	    product_id,
	    warehouse_place_id,
	    type,
	    quantity,
	    reference_id,
	    reference_type,
	    reason,
	    created_at
	FROM stock_movement
	WHERE product_id = ?
	ORDER BY created_at DESC
	`)

	mock.ExpectQuery(query).
		WithArgs(1).
		WillReturnRows(rows)

	result, err := dao.GetMovementsByProduct(1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), result[0].ProductID)
	assert.Equal(t, 10, result[0].Quantity)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStockDAO_GetAllMovements(t *testing.T) {

	db, mock, dao := newMock()
	defer db.Close()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"product_id",
		"warehouse_place_id",
		"type",
		"quantity",
		"reference_id",
		"reference_type",
		"reason",
		"created_at",
	}).
		AddRow(1, 1, "IN", 10, nil, nil, "purchase", now).
		AddRow(2, 1, "OUT", 3, nil, nil, "sale", now)

	query := regexp.QuoteMeta(`
	SELECT 
	    product_id,
	    warehouse_place_id,
	    type,
	    quantity,
	    reference_id,
	    reference_type,
	    reason,
	    created_at
	FROM stock_movement
	ORDER BY created_at DESC
	`)

	mock.ExpectQuery(query).
		WillReturnRows(rows)

	result, err := dao.GetAllMovements()

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	assert.Equal(t, int64(1), result[0].ProductID)
	assert.Equal(t, int64(2), result[1].ProductID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStockDAO_GetMovementsByProductAndWarehouse(t *testing.T) {

	db, mock, dao := newMock()
	defer db.Close()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"product_id",
		"warehouse_place_id",
		"type",
		"quantity",
		"reference_id",
		"reference_type",
		"reason",
		"created_at",
	}).AddRow(
		1, 2, "IN", 5, nil, nil, "adjustment", now,
	)

	query := regexp.QuoteMeta(`
	SELECT 
	    product_id,
	    warehouse_place_id,
	    type,
	    quantity,
	    reference_id,
	    reference_type,
	    reason,
	    created_at
	FROM stock_movement
	WHERE product_id = ? 
	  AND warehouse_place_id = ?
	ORDER BY created_at DESC
	`)

	mock.ExpectQuery(query).
		WithArgs(1, 2).
		WillReturnRows(rows)

	result, err := dao.GetMovementsByProductAndWarehouse(1, 2)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 5, result[0].Quantity)

	assert.NoError(t, mock.ExpectationsWereMet())
}
