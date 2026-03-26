package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/stretchr/testify/assert"
)

func newStockDAO(db *sql.DB) StockDAO {
	return &stockDAO{db: db}
}

func TestStock_Upsert(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newStockDAO(db)

	mock.ExpectExec("INSERT INTO stock").
		WithArgs(1, 1, 10, true, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	stock := &model.Stock{
		ProductID:   1,
		WarehouseID: 1,
		Quantity:    10,
		Active:      true,
	}

	err := dao.Upsert(stock)

	assert.NoError(t, err)
}

func TestStock_Get(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newStockDAO(db)

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"product_id",
		"warehouse_id",
		"quantity",
		"active",
		"updated_at",
	}).AddRow(1, 1, 20, true, now)

	mock.ExpectQuery("SELECT product_id").
		WithArgs(1, 1).
		WillReturnRows(rows)

	result, err := dao.Get(1, 1)

	assert.NoError(t, err)
	assert.Equal(t, 20, result.Quantity)
}

func TestStock_List(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newStockDAO(db)

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"product_id",
		"warehouse_id",
		"quantity",
		"active",
		"updated_at",
	}).AddRow(1, 1, 5, true, now)

	mock.ExpectQuery("SELECT product_id").
		WillReturnRows(rows)

	list, err := dao.List()

	assert.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestStock_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newStockDAO(db)

	mock.ExpectExec("DELETE FROM stock").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Delete(1, 1)

	assert.NoError(t, err)
}

func newStockMovementDAO(db *sql.DB) StockMovementDAO {
	return &stockMovementDAO{db: db}
}

func TestStockMovement_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newStockMovementDAO(db)

	now := time.Now()

	mock.ExpectExec("INSERT INTO stock_movement").
		WithArgs(1, 1, model.StockMovementIn, 10, now, "entrada").
		WillReturnResult(sqlmock.NewResult(1, 1))

	m := &model.StockMovement{
		ProductID:   1,
		WarehouseID: 1,
		Type:        model.StockMovementIn,
		Quantity:    10,
		CreatedAt:   now,
		Reason:      "entrada",
	}

	err := dao.Create(m)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), m.ID)
}

func TestStockMovement_ListByProduct(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newStockMovementDAO(db)

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id",
		"product_id",
		"warehouse_id",
		"type",
		"quantity",
		"created_at",
		"reason",
	}).AddRow(1, 1, 1, "IN", 10, now, "entrada")

	mock.ExpectQuery("SELECT id, product_id").
		WithArgs(1).
		WillReturnRows(rows)

	list, err := dao.ListByProduct(1)

	assert.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestStockMovement_ListByWarehouse(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newStockMovementDAO(db)

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id",
		"product_id",
		"warehouse_id",
		"type",
		"quantity",
		"created_at",
		"reason",
	}).AddRow(1, 1, 1, "OUT", 5, now, "saida")

	mock.ExpectQuery("SELECT id, product_id").
		WithArgs(1).
		WillReturnRows(rows)

	list, err := dao.ListByWarehouse(1)

	assert.NoError(t, err)
	assert.Len(t, list, 1)
}
