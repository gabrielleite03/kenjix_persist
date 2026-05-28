package repository

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
)

type SalesOrderInvoiceRepository interface {
	Create(invoice *model.SalesOrderInvoice) (*model.SalesOrderInvoice, error)
	Update(invoice *model.SalesOrderInvoice) error
}

type salesOrderInvoiceRepository struct {
	db *sql.DB
}

func NewSalesOrderInvoiceRepository() SalesOrderInvoiceRepository {
	return &salesOrderInvoiceRepository{db: config.NewDatabaseConfig().DB}
}

func (r *salesOrderInvoiceRepository) Create(invoice *model.SalesOrderInvoice) (*model.SalesOrderInvoice, error) {
	now := time.Now()
	invoice.CreatedAt = now
	invoice.UpdatedAt = now

	query := `
		INSERT INTO sales_order_invoice (
			sales_order_id,
			nfe_key,
			number,
			series,
			status,
			status_code,
			status_reason,
			protocol,
			issued_at,
			authorized_at,
			cancelled_at,
			xml_path,
			pdf_path,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		invoice.SalesOrderID,
		invoice.NFeKey,
		invoice.Number,
		invoice.Series,
		invoice.Status,
		invoice.StatusCode,
		invoice.StatusReason,
		invoice.Protocol,
		invoice.IssuedAt,
		invoice.AuthorizedAt,
		invoice.CancelledAt,
		invoice.XMLPath,
		invoice.PDFPath,
		invoice.CreatedAt,
		invoice.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	invoice.ID, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (r *salesOrderInvoiceRepository) Update(invoice *model.SalesOrderInvoice) error {
	invoice.UpdatedAt = time.Now()

	query := `
		UPDATE sales_order_invoice
		SET
			nfe_key = ?,
			number = ?,
			series = ?,
			status = ?,
			status_code = ?,
			status_reason = ?,
			protocol = ?,
			issued_at = ?,
			authorized_at = ?,
			cancelled_at = ?,
			xml_path = ?,
			pdf_path = ?,
			updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		invoice.NFeKey,
		invoice.Number,
		invoice.Series,
		invoice.Status,
		invoice.StatusCode,
		invoice.StatusReason,
		invoice.Protocol,
		invoice.IssuedAt,
		invoice.AuthorizedAt,
		invoice.CancelledAt,
		invoice.XMLPath,
		invoice.PDFPath,
		invoice.UpdatedAt,
		invoice.ID,
	)

	return err
}

func getenv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}
