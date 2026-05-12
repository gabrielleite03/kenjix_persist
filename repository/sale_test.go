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

	customerName := "Gabriel Leite"
	customerDocument := "12345678900"

	marketplaceID := int64(1)

	externalOrderID := "2000016386215972"
	externalPackID := "PACK123"

	paymentStatus := "paid"
	deliveryStatus := "no_shipping"

	mock.ExpectExec("INSERT INTO sales_order").
		WithArgs(
			"100.00",
			"10.00",
			"OPEN",
			int64(1),
			true,

			&customerName,
			&customerDocument,

			&marketplaceID,
			&externalOrderID,
			&externalPackID,

			&paymentStatus,
			&deliveryStatus,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s := &model.SalesOrder{
		Price:           decimal.NewFromFloat(100),
		Discount:        decimal.NewFromFloat(10),
		Status:          "OPEN",
		PaymentMethodID: 1,
		Active:          true,

		CustomerName:     &customerName,
		CustomerDocument: &customerDocument,

		MarketplaceID:   &marketplaceID,
		ExternalOrderID: &externalOrderID,
		ExternalPackID:  &externalPackID,

		PaymentStatus:  &paymentStatus,
		DeliveryStatus: &deliveryStatus,
	}

	created, err := dao.Create(s)

	assert.NoError(t, err)
	assert.NotNil(t, created)

	assert.Equal(t, int64(1), created.ID)
	assert.Equal(t, int64(1), s.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSalesOrder_GetByID(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newSalesDAO(db)

	now := time.Now()

	updatedAt := now

	customerName := "Gabriel Leite"
	customerDocument := "12345678900"

	marketplaceID := int64(1)

	externalOrderID := "2000016386215972"
	externalPackID := "PACK123"

	paymentStatus := "paid"
	deliveryStatus := "no_shipping"

	rows := sqlmock.NewRows([]string{
		"id",
		"price",
		"discount",
		"status",
		"payment_method_id",
		"active",
		"created_at",
		"customer_name",
		"customer_document",
		"marketplace_id",
		"external_order_id",
		"external_pack_id",
		"payment_status",
		"delivery_status",
		"updated_at",

		"pm_id",
		"pm_name",
		"pm_active",

		"m_id",
		"m_name",
		"m_description",
		"m_commission_rate",
		"m_active",
	}).AddRow(

		1,
		"100.00",
		"5.00",
		"OPEN",
		1,
		true,
		now,

		customerName,
		customerDocument,

		marketplaceID,
		externalOrderID,
		externalPackID,

		paymentStatus,
		deliveryStatus,

		updatedAt,

		1,
		"PIX",
		true,

		1,
		"Mercado Livre",
		"Marketplace ML",
		"16.50",
		true,
	)

	mock.ExpectQuery("SELECT(.+)FROM sales_order").
		WithArgs(1).
		WillReturnRows(rows)

	mock.ExpectQuery("FROM sales_order_invoice").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	result, err := dao.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, int64(1), result.ID)

	assert.True(t,
		result.Price.Equal(decimal.NewFromFloat(100)),
	)

	assert.True(t,
		result.Discount.Equal(decimal.NewFromFloat(5)),
	)

	assert.Equal(t, "OPEN", result.Status)

	assert.NotNil(t, result.CustomerName)
	assert.Equal(t, customerName, *result.CustomerName)

	assert.NotNil(t, result.CustomerDocument)
	assert.Equal(t, customerDocument, *result.CustomerDocument)

	assert.NotNil(t, result.MarketplaceID)
	assert.Equal(t, marketplaceID, *result.MarketplaceID)

	assert.NotNil(t, result.ExternalOrderID)
	assert.Equal(t, externalOrderID, *result.ExternalOrderID)

	assert.NotNil(t, result.ExternalPackID)
	assert.Equal(t, externalPackID, *result.ExternalPackID)

	assert.NotNil(t, result.PaymentStatus)
	assert.Equal(t, paymentStatus, *result.PaymentStatus)

	assert.NotNil(t, result.DeliveryStatus)
	assert.Equal(t, deliveryStatus, *result.DeliveryStatus)

	assert.NotNil(t, result.PaymentMethod)
	assert.Equal(t, int64(1), result.PaymentMethod.ID)
	assert.Equal(t, "PIX", result.PaymentMethod.Name)

	assert.NotNil(t, result.Marketplace)
	assert.Equal(t, int64(1), result.Marketplace.ID)
	assert.Equal(t, "Mercado Livre", result.Marketplace.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

//
// SalesOrderItem Tests
//

func TestSalesOrderItem_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newSalesItemDAO(db)

	mock.ExpectExec("INSERT INTO sales_order_item").
		WithArgs(
			int64(1),
			int64(10),
			2,
			"50.00",
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	item := &model.SalesOrderItem{
		SalesOrderID:   1,
		PurchaseItemID: 10,
		Quantity:       2,
		UnitPrice:      decimal.NewFromFloat(50),
	}

	err := dao.Create(item)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSalesOrderItem_ListByOrder(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newSalesItemDAO(db)

	rows := sqlmock.NewRows([]string{
		"sales_order_id",
		"purchase_item_id",
		"quantity",
		"unit_price",

		"id",
		"purchase_id",
		"product_id",
		"purchase_quantity",
		"cost_price",
		"total",
		"cost_center_id",
	}).AddRow(
		1,
		10,
		2,
		"50.00",

		10,
		1,
		99,
		"5.0000",
		"25.0000",
		"125.00",
		nil,
	)

	mock.ExpectQuery("FROM sales_order_item").
		WithArgs(1).
		WillReturnRows(rows)

	list, err := dao.ListByOrder(1)

	assert.NoError(t, err)

	assert.Len(t, list, 1)

	assert.Equal(t, int64(1), list[0].SalesOrderID)
	assert.Equal(t, int64(10), list[0].PurchaseItemID)

	assert.Equal(t, 2, list[0].Quantity)

	assert.True(
		t,
		list[0].UnitPrice.Equal(decimal.NewFromFloat(50)),
	)

	assert.NotNil(t, list[0].PurchaseItem)

	assert.Equal(
		t,
		int64(10),
		list[0].PurchaseItem.ID,
	)

	assert.Equal(
		t,
		int64(99),
		list[0].PurchaseItem.ProductID,
	)

	assert.True(
		t,
		list[0].PurchaseItem.Quantity.Equal(
			decimal.RequireFromString("5.0000"),
		),
	)

	assert.True(
		t,
		list[0].PurchaseItem.CostPrice.Equal(
			decimal.RequireFromString("25.0000"),
		),
	)

	assert.True(
		t,
		list[0].PurchaseItem.Total.Equal(
			decimal.RequireFromString("125.00"),
		),
	)

	assert.Nil(t, list[0].PurchaseItem.CostCenterID)

	assert.NoError(t, mock.ExpectationsWereMet())
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
