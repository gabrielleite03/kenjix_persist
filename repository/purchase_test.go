package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestPurchaseDAO_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewPurchaseDAOWithDB(db)

	p := &model.Purchase{
		InvoiceType: "NF",
		SupplierID:  1,
		Status:      "OPEN",
		Total:       decimal.NewFromFloat(100),
		Items: []model.PurchaseItem{
			{
				ProductID: 1,
				Quantity:  decimal.NewFromInt(2),
				CostPrice: decimal.NewFromFloat(50),
				Total:     decimal.NewFromFloat(100),
			},
		},
	}

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO purchase
		(invoice_number, invoice_type, supplier_id, status, total)
		VALUES (?, ?, ?, ?, ?)
	`)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta(`
			INSERT INTO purchase_item
			(purchase_id, product_id, quantity, cost_price, total, cost_center_id)
			VALUES (?, ?, ?, ?, ?, ?)
		`)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = dao.Create(p)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), p.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPurchaseDAO_FindByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := NewPurchaseDAOWithDB(db)

	rows := sqlmock.NewRows([]string{
		"id", "invoice_number", "invoice_type", "supplier_id", "status", "total", "created_at",
	}).AddRow(
		1, nil, "NF", 1, "OPEN", 100.0, time.Now(),
	)

	mock.ExpectQuery("SELECT id").
		WithArgs(1).
		WillReturnRows(rows)

	itemRows := sqlmock.NewRows([]string{
		"id", "purchase_id", "product_id", "quantity", "cost_price", "total", "cost_center_id",
	}).AddRow(
		1, 1, 1, 2.0, 50.0, 100.0, nil,
	)

	mock.ExpectQuery("FROM purchase_item").
		WithArgs(1).
		WillReturnRows(itemRows)

	result, err := dao.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Len(t, result.Items, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPurchaseDAO_FindAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := NewPurchaseDAOWithDB(db)

	rows := sqlmock.NewRows([]string{
		"id", "invoice_number", "invoice_type", "supplier_id", "status", "total", "created_at",
	}).AddRow(
		1, nil, "NF", 1, "OPEN", 100.0, time.Now(),
	)

	mock.ExpectQuery("SELECT id").
		WillReturnRows(rows)

	list, err := dao.FindAll()

	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}
