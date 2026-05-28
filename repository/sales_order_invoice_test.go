package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/stretchr/testify/assert"
)

func newSalesOrderInvoiceRepository(db *sql.DB) SalesOrderInvoiceRepository {
	return &salesOrderInvoiceRepository{db: db}
}

func TestSalesOrderInvoice_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := newSalesOrderInvoiceRepository(db)

	nfeKey := "35240100000000000000550010000000011000000010"
	number := int64(123)
	series := 1
	statusCode := "100"
	statusReason := "Autorizado"
	protocol := "135240000000000"
	issuedAt := time.Now()
	authorizedAt := issuedAt.Add(time.Minute)
	xmlPath := "/tmp/nfe.xml"
	pdfPath := "/tmp/nfe.pdf"

	mock.ExpectExec("INSERT INTO sales_order_invoice").
		WithArgs(
			int64(1),
			&nfeKey,
			&number,
			&series,
			"AUTHORIZED",
			&statusCode,
			&statusReason,
			&protocol,
			&issuedAt,
			&authorizedAt,
			nil,
			&xmlPath,
			&pdfPath,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(10, 1))

	invoice := &model.SalesOrderInvoice{
		SalesOrderID: int64(1),
		NFeKey:       &nfeKey,
		Number:       &number,
		Series:       &series,
		Status:       "AUTHORIZED",
		StatusCode:   &statusCode,
		StatusReason: &statusReason,
		Protocol:     &protocol,
		IssuedAt:     &issuedAt,
		AuthorizedAt: &authorizedAt,
		CancelledAt:  nil,
		XMLPath:      &xmlPath,
		PDFPath:      &pdfPath,
	}

	created, err := repo.Create(invoice)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, int64(10), created.ID)
	assert.False(t, created.CreatedAt.IsZero())
	assert.False(t, created.UpdatedAt.IsZero())
	assert.Equal(t, created.CreatedAt, created.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSalesOrderInvoice_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := newSalesOrderInvoiceRepository(db)

	nfeKey := "35240100000000000000550010000000011000000010"
	number := int64(123)
	series := 1
	statusCode := "100"
	statusReason := "Autorizado"
	protocol := "135240000000000"
	issuedAt := time.Now()
	authorizedAt := issuedAt.Add(time.Minute)
	cancelledAt := issuedAt.Add(time.Hour)
	xmlPath := "/tmp/nfe.xml"
	pdfPath := "/tmp/nfe.pdf"

	mock.ExpectExec("UPDATE sales_order_invoice").
		WithArgs(
			&nfeKey,
			&number,
			&series,
			"CANCELLED",
			&statusCode,
			&statusReason,
			&protocol,
			&issuedAt,
			&authorizedAt,
			&cancelledAt,
			&xmlPath,
			&pdfPath,
			sqlmock.AnyArg(),
			int64(10),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	invoice := &model.SalesOrderInvoice{
		ID:           int64(10),
		SalesOrderID: int64(1),
		NFeKey:       &nfeKey,
		Number:       &number,
		Series:       &series,
		Status:       "CANCELLED",
		StatusCode:   &statusCode,
		StatusReason: &statusReason,
		Protocol:     &protocol,
		IssuedAt:     &issuedAt,
		AuthorizedAt: &authorizedAt,
		CancelledAt:  &cancelledAt,
		XMLPath:      &xmlPath,
		PDFPath:      &pdfPath,
		UpdatedAt:    time.Now().Add(-time.Hour),
	}

	oldUpdatedAt := invoice.UpdatedAt

	err := repo.Update(invoice)

	assert.NoError(t, err)
	assert.True(t, invoice.UpdatedAt.After(oldUpdatedAt))
	assert.NoError(t, mock.ExpectationsWereMet())
}
