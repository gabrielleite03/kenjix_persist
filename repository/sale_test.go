package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

//
// helpers
//

func newPaymentDAO(db *sql.DB) PaymentMethodDAO {
	return &paymentMethodDAO{db: db}
}

func newSalesDAO(db *sql.DB) SalesOrderDAO {
	return &salesOrderDAO{db: db}
}

func newSalesItemDAO(db *sql.DB) SalesOrderItemDAO {
	return &salesOrderItemDAO{db: db}
}

//
// PaymentMethod Tests
//

func TestPaymentMethod_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newPaymentDAO(db)

	mock.ExpectExec("INSERT INTO payment_method").
		WithArgs("PIX").
		WillReturnResult(sqlmock.NewResult(1, 1))

	p := &model.PaymentMethod{Name: "PIX"}

	err := dao.Create(p)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), p.ID)
}

func TestPaymentMethod_List(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newPaymentDAO(db)

	rows := sqlmock.NewRows([]string{"id", "name", "active"}).
		AddRow(1, "PIX", true).
		AddRow(2, "CREDIT", true)

	mock.ExpectQuery("SELECT id, name, active FROM payment_method").
		WillReturnRows(rows)

	list, err := dao.List()

	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

//
// SalesOrder Tests
//

func TestSalesOrder_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newSalesDAO(db)

	mock.ExpectExec("INSERT INTO sales_order").
		WithArgs(
			"100.00",
			"10.00",
			"OPEN",
			1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s := &model.SalesOrder{
		Price:           decimal.NewFromFloat(100),
		Discount:        decimal.NewFromFloat(10),
		Status:          "OPEN",
		PaymentMethodID: 1,
	}

	err := dao.Create(s)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), s.ID)
}

func TestSalesOrder_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newSalesDAO(db)

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "price", "discount", "status",
		"payment_method_id", "active", "created_at",
	}).AddRow(
		1,
		"100.00",
		"5.00",
		"OPEN",
		1,
		true,
		now,
	)

	mock.ExpectQuery("FROM sales_order").
		WithArgs(1).
		WillReturnRows(rows)

	result, err := dao.GetByID(1)

	assert.NoError(t, err)
	assert.True(t, result.Price.Equal(decimal.NewFromFloat(100)))
}

//
// SalesOrderItem Tests
//

func TestSalesOrderItem_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newSalesItemDAO(db)

	mock.ExpectExec("INSERT INTO sales_order_item").
		WithArgs(1, 10, 2, "50.00").
		WillReturnResult(sqlmock.NewResult(1, 1))

	item := &model.SalesOrderItem{
		SalesOrderID: 1,
		ProductID:    10,
		Quantity:     2,
		UnitPrice:    decimal.NewFromFloat(50),
	}

	err := dao.Create(item)

	assert.NoError(t, err)
}

func TestSalesOrderItem_ListByOrder(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newSalesItemDAO(db)

	rows := sqlmock.NewRows([]string{
		"sales_order_id",
		"product_id",
		"quantity",
		"unit_price",
	}).AddRow(1, 10, 2, "50.00")

	mock.ExpectQuery("FROM sales_order_item").
		WithArgs(1).
		WillReturnRows(rows)

	list, err := dao.ListByOrder(1)

	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.True(t, list[0].UnitPrice.Equal(decimal.NewFromFloat(50)))
}

func TestSalesOrderItem_DeleteByOrder(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newSalesItemDAO(db)

	mock.ExpectExec("DELETE FROM sales_order_item").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.DeleteByOrder(1)

	assert.NoError(t, err)
}
