package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/stretchr/testify/assert"
)

func setupMock() (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, _ := sqlmock.New()
	return db, mock, func() { db.Close() }
}

func TestStockDAO_Create(t *testing.T) {
	db, mock, close := setupMock()
	defer close()

	dao := &StockDAO{db: db}

	stock := &model.Stock{
		Product:        model.Product{ID: 1},
		WarehousePlace: model.WarehousePlace{ID: 2},
		PurchaseItem:   model.PurchaseItem{ID: 3},
		Quantity:       10,
	}

	mock.ExpectExec("INSERT INTO stock").
		WithArgs(1, 2, 3, 10).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := dao.Create(stock)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStockDAO_GetByID(t *testing.T) {
	db, mock, close := setupMock()
	defer close()

	dao := &StockDAO{db: db}

	rows := sqlmock.NewRows([]string{
		"id", "product_id", "warehouse_place_id",
		"purchase_item_id", "quantity", "active", "updated_at",
	}).AddRow(1, 10, 20, 30, 5, true, time.Now())

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnRows(rows)

	result, err := dao.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, 5, result.Quantity)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStockDAO_GetAllActive(t *testing.T) {
	db, mock, close := setupMock()
	defer close()

	dao := &StockDAO{db: db}

	rows := sqlmock.NewRows([]string{
		"id", "product_id", "warehouse_place_id",
		"purchase_item_id", "quantity", "active", "updated_at",
	}).
		AddRow(1, 1, 1, 1, 10, true, time.Now()).
		AddRow(2, 1, 2, 2, 20, true, time.Now())

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	list, err := dao.GetAllActive()

	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, 10, list[0].Quantity)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStockDAO_GetGroupedByProductAndWarehouse(t *testing.T) {
	db, mock, close := setupMock()
	defer close()

	dao := &StockDAO{db: db}

	rows := sqlmock.NewRows([]string{
		"product_id", "warehouse_place_id", "quantity",
	}).
		AddRow(1, 1, 30).
		AddRow(1, 2, 20)

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	list, err := dao.GetGroupedByProductAndWarehouse()

	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, 30, list[0].Quantity)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStockMovementDAO_Create(t *testing.T) {
	db, mock, close := setupMock()
	defer close()

	dao := &StockMovementDAO{db: db}

	refID := int64(99)
	refType := "PURCHASE"

	m := &model.StockMovement{
		ProductID:        1,
		WarehousePlaceID: 2,
		PurchseItemID:    3,
		Type:             "IN",
		Quantity:         10,
		ReferenceID:      &refID,
		ReferenceType:    &refType,
		Reason:           "Compra",
	}

	mock.ExpectExec("INSERT INTO stock_movement").
		WithArgs(1, 2, 3, "IN", 10, &refID, &refType, "Compra").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := dao.Create(m)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStockMovementDAO_GetByProduct(t *testing.T) {
	db, mock, close := setupMock()
	defer close()

	dao := &StockMovementDAO{db: db}

	rows := sqlmock.NewRows([]string{
		"id", "product_id", "warehouse_place_id",
		"purchase_item_id", "type", "quantity",
		"reference_id", "reference_type", "reason", "created_at",
	}).AddRow(1, 1, 1, 1, "IN", 10, nil, nil, "", time.Now())

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnRows(rows)

	list, err := dao.GetByProduct(1)

	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, 10, list[0].Quantity)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAllEager(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := &StockMovementDAO{db: db}
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "type", "quantity", "reference_id", "reference_type", "reason", "created_at",

		"id", "name", "sku", "price", "marca", "description", "active", "volume", "category_id",

		"id", "name", "active", "warehouse_place_type_id", "warehouse_id", "capacity",

		"id", "purchase_id", "product_id", "quantity", "cost_price", "total", "cost_center_id",
	}).AddRow(
		1, "IN", 10, nil, nil, "entrada", now,

		// product
		100, "Produto A", "SKU123", "10.50", "MarcaX", "Desc", true, "1.5", nil,

		// warehouse_place
		200, "Prateleira 1", true, nil, nil, nil,

		// purchase_item
		300, 400, 100, "10", "5.00", "50.00", nil,
	)

	mock.ExpectQuery("SELECT (.+) FROM stock_movement").
		WillReturnRows(rows)

	result, err := dao.FindAllEager()

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	sm := result[0]

	assert.Equal(t, int64(1), sm.ID)
	assert.Equal(t, model.StockMovementType("IN"), sm.Type)
	assert.Equal(t, 10, sm.Quantity)

	assert.Equal(t, int64(100), sm.Product.ID)
	assert.Equal(t, "Produto A", sm.Product.Name)

	assert.Equal(t, int64(200), sm.WarehousePlace.ID)
	assert.Equal(t, "Prateleira 1", sm.WarehousePlace.Name)

	assert.Equal(t, int64(300), sm.PurchaseItem.ID)
	assert.Equal(t, int64(400), sm.PurchaseItem.PurchaseID)

	assert.NoError(t, mock.ExpectationsWereMet())
}
