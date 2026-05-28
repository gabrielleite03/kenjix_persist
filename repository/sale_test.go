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

func TestPaymentMethod_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newPaymentDAO(db)

	rows := sqlmock.NewRows([]string{"id", "name", "active"}).
		AddRow(1, "PIX", true)

	mock.ExpectQuery("SELECT id, name, active FROM payment_method WHERE id=\\?").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	result, err := dao.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "PIX", result.Name)
	assert.True(t, result.Active)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPaymentMethod_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newPaymentDAO(db)

	mock.ExpectExec("UPDATE payment_method SET name=\\?, active=\\? WHERE id=\\?").
		WithArgs("CREDIT", false, int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Update(&model.PaymentMethod{
		ID:     1,
		Name:   "CREDIT",
		Active: false,
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPaymentMethod_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newPaymentDAO(db)

	mock.ExpectExec("DELETE FROM payment_method WHERE id=\\?").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Delete(1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
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

func TestSalesOrder_GetByID_WithInvoice(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newSalesDAO(db)

	now := time.Now()
	nfeKey := "35240100000000000000550010000000011000000010"
	number := int64(123)
	series := int64(1)
	statusCode := "100"
	statusReason := "Autorizado"
	protocol := "135240000000000"
	xmlPath := "/tmp/nfe.xml"
	pdfPath := "/tmp/nfe.pdf"

	saleRows := sqlmock.NewRows([]string{
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
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		1,
		"PIX",
		true,
		1,
		"Mercado Livre",
		"Marketplace ML",
		"16.50",
		true,
	)

	invoiceRows := sqlmock.NewRows([]string{
		"id",
		"sales_order_id",
		"nfe_key",
		"number",
		"series",
		"status",
		"status_code",
		"status_reason",
		"protocol",
		"issued_at",
		"authorized_at",
		"cancelled_at",
		"xml_path",
		"pdf_path",
		"created_at",
		"updated_at",
	}).AddRow(
		10,
		1,
		nfeKey,
		number,
		series,
		"AUTHORIZED",
		statusCode,
		statusReason,
		protocol,
		now,
		now,
		nil,
		xmlPath,
		pdfPath,
		now,
		now,
	)

	mock.ExpectQuery("SELECT(.+)FROM sales_order").
		WithArgs(1).
		WillReturnRows(saleRows)

	mock.ExpectQuery("FROM sales_order_invoice").
		WithArgs(int64(1)).
		WillReturnRows(invoiceRows)

	result, err := dao.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Invoice)
	assert.Equal(t, int64(10), result.Invoice.ID)
	assert.Equal(t, int64(1), result.Invoice.SalesOrderID)
	assert.Equal(t, nfeKey, *result.Invoice.NFeKey)
	assert.Equal(t, number, *result.Invoice.Number)
	assert.Equal(t, int(series), *result.Invoice.Series)
	assert.Equal(t, "AUTHORIZED", result.Invoice.Status)
	assert.Equal(t, statusCode, *result.Invoice.StatusCode)
	assert.Equal(t, statusReason, *result.Invoice.StatusReason)
	assert.Equal(t, protocol, *result.Invoice.Protocol)
	assert.NotNil(t, result.Invoice.IssuedAt)
	assert.NotNil(t, result.Invoice.AuthorizedAt)
	assert.Nil(t, result.Invoice.CancelledAt)
	assert.Equal(t, xmlPath, *result.Invoice.XMLPath)
	assert.Equal(t, pdfPath, *result.Invoice.PDFPath)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSalesOrder_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newSalesDAO(db)

	customerName := "Gabriel Leite"
	customerDocument := "12345678900"
	marketplaceID := int64(1)
	externalOrderID := "2000016386215972"
	externalPackID := "PACK123"
	paymentStatus := "paid"
	deliveryStatus := "delivered"

	mock.ExpectExec("UPDATE sales_order").
		WithArgs(
			"150.00",
			"15.00",
			"CLOSED",
			int64(2),
			false,
			&customerName,
			&customerDocument,
			&marketplaceID,
			&externalOrderID,
			&externalPackID,
			&paymentStatus,
			&deliveryStatus,
			int64(1),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Update(&model.SalesOrder{
		ID:               1,
		Price:            decimal.NewFromFloat(150),
		Discount:         decimal.NewFromFloat(15),
		Status:           "CLOSED",
		PaymentMethodID:  2,
		Active:           false,
		CustomerName:     &customerName,
		CustomerDocument: &customerDocument,
		MarketplaceID:    &marketplaceID,
		ExternalOrderID:  &externalOrderID,
		ExternalPackID:   &externalPackID,
		PaymentStatus:    &paymentStatus,
		DeliveryStatus:   &deliveryStatus,
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSalesOrder_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newSalesDAO(db)

	mock.ExpectExec("DELETE FROM sales_order WHERE id=\\?").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Delete(1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSalesOrder_List(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newSalesDAO(db)

	now := time.Now()
	deletedAt := now.Add(time.Hour)
	logo := "logo.png"
	apiURL := "https://api.example.com"
	apiKey := "api-key"
	apiSecret := "api-secret"
	apiEndpoint := "/orders"

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
		"m_commission_rate",
		"m_logo",
		"m_status",
		"m_integration_type",
		"m_api_url",
		"m_api_key",
		"m_api_secret",
		"m_api_endpoint",
		"m_created_at",
		"m_deleted_at",
	}).AddRow(
		1,
		"200.00",
		"20.00",
		"OPEN",
		1,
		true,
		now,
		"Gabriel Leite",
		"12345678900",
		1,
		"2000016386215972",
		"PACK123",
		"paid",
		"delivered",
		now,
		1,
		"PIX",
		true,
		1,
		"Mercado Livre",
		"16.50",
		logo,
		"ACTIVE",
		"API",
		apiURL,
		apiKey,
		apiSecret,
		apiEndpoint,
		now,
		deletedAt,
	)

	mock.ExpectQuery("SELECT(.+)FROM sales_order").
		WillReturnRows(rows)

	mock.ExpectQuery("FROM sales_order_invoice").
		WithArgs(int64(1)).
		WillReturnError(sql.ErrNoRows)

	list, err := dao.List()

	assert.NoError(t, err)
	assert.Len(t, list, 1)

	result := list[0]
	assert.Equal(t, int64(1), result.ID)
	assert.True(t, result.Price.Equal(decimal.NewFromFloat(200)))
	assert.True(t, result.Discount.Equal(decimal.NewFromFloat(20)))
	assert.Equal(t, "OPEN", result.Status)
	assert.Equal(t, int64(1), result.PaymentMethodID)
	assert.True(t, result.Active)
	assert.Equal(t, "Gabriel Leite", *result.CustomerName)
	assert.Equal(t, "12345678900", *result.CustomerDocument)
	assert.Equal(t, int64(1), *result.MarketplaceID)
	assert.Equal(t, "2000016386215972", *result.ExternalOrderID)
	assert.Equal(t, "PACK123", *result.ExternalPackID)
	assert.Equal(t, "paid", *result.PaymentStatus)
	assert.Equal(t, "delivered", *result.DeliveryStatus)
	assert.NotNil(t, result.UpdatedAt)

	assert.NotNil(t, result.PaymentMethod)
	assert.Equal(t, int64(1), result.PaymentMethod.ID)
	assert.Equal(t, "PIX", result.PaymentMethod.Name)

	assert.NotNil(t, result.Marketplace)
	assert.Equal(t, int64(1), result.Marketplace.ID)
	assert.Equal(t, "Mercado Livre", result.Marketplace.Name)
	assert.True(t, result.Marketplace.CommissionRate.Equal(decimal.RequireFromString("16.50")))
	assert.Equal(t, logo, *result.Marketplace.Logo)
	assert.Equal(t, "ACTIVE", result.Marketplace.Status)
	assert.Equal(t, "API", result.Marketplace.IntegrationType)
	assert.Equal(t, apiURL, *result.Marketplace.APIURL)
	assert.Equal(t, apiKey, *result.Marketplace.APIKey)
	assert.Equal(t, apiSecret, *result.Marketplace.APISecret)
	assert.Equal(t, apiEndpoint, *result.Marketplace.APIEndpoint)
	assert.NotNil(t, result.Marketplace.DeletedAt)
	assert.Nil(t, result.Invoice)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSalesOrder_GetInvoiceBySalesOrderID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &salesOrderDAO{db: db}

	now := time.Now()
	statusCode := "100"

	rows := sqlmock.NewRows([]string{
		"id",
		"sales_order_id",
		"nfe_key",
		"number",
		"series",
		"status",
		"status_code",
		"status_reason",
		"protocol",
		"issued_at",
		"authorized_at",
		"cancelled_at",
		"xml_path",
		"pdf_path",
		"created_at",
		"updated_at",
	}).AddRow(
		10,
		1,
		nil,
		nil,
		nil,
		"AUTHORIZED",
		statusCode,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		now,
		now,
	)

	mock.ExpectQuery("FROM sales_order_invoice").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	result, err := dao.GetInvoiceBySalesOrderID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(10), result.ID)
	assert.Equal(t, int64(1), result.SalesOrderID)
	assert.Equal(t, "AUTHORIZED", result.Status)
	assert.Equal(t, statusCode, *result.StatusCode)
	assert.Nil(t, result.NFeKey)
	assert.Nil(t, result.Number)
	assert.Nil(t, result.Series)
	assert.Nil(t, result.StatusReason)
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

func TestSalesOrderItem_ListByOrder_WithCostCenter(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newSalesItemDAO(db)

	costCenterID := int64(7)

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
		costCenterID,
	)

	mock.ExpectQuery("FROM sales_order_item").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	list, err := dao.ListByOrder(1)

	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.NotNil(t, list[0].PurchaseItem)
	assert.NotNil(t, list[0].PurchaseItem.CostCenterID)
	assert.Equal(t, costCenterID, *list[0].PurchaseItem.CostCenterID)
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

func TestSalesOrderItem_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := newSalesItemDAO(db)

	now := time.Now()
	customerName := "Gabriel Leite"
	customerDocument := "12345678900"
	marketplaceID := int64(1)
	externalOrderID := "2000016386215972"
	externalPackID := "PACK123"
	paymentStatus := "paid"
	deliveryStatus := "delivered"

	rows := sqlmock.NewRows([]string{
		"id",
		"price",
		"discount",
		"status",
		"payment_method_id",
		"active",
		"customer_name",
		"customer_document",
		"marketplace_id",
		"external_order_id",
		"external_pack_id",
		"payment_status",
		"delivery_status",
		"created_at",
		"updated_at",
	}).AddRow(
		1,
		"100.00",
		"10.00",
		"OPEN",
		1,
		true,
		customerName,
		customerDocument,
		marketplaceID,
		externalOrderID,
		externalPackID,
		paymentStatus,
		deliveryStatus,
		now,
		now,
	)

	mock.ExpectQuery("FROM sales_order").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	result, err := dao.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.True(t, result.Price.Equal(decimal.NewFromFloat(100)))
	assert.True(t, result.Discount.Equal(decimal.NewFromFloat(10)))
	assert.Equal(t, "OPEN", result.Status)
	assert.Equal(t, int64(1), result.PaymentMethodID)
	assert.True(t, result.Active)
	assert.Equal(t, customerName, *result.CustomerName)
	assert.Equal(t, customerDocument, *result.CustomerDocument)
	assert.Equal(t, marketplaceID, *result.MarketplaceID)
	assert.Equal(t, externalOrderID, *result.ExternalOrderID)
	assert.Equal(t, externalPackID, *result.ExternalPackID)
	assert.Equal(t, paymentStatus, *result.PaymentStatus)
	assert.Equal(t, deliveryStatus, *result.DeliveryStatus)
	assert.NotNil(t, result.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}
