package repository

import (
	"database/sql"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
)

type PurchaseDAO struct {
	db *sql.DB
}

func NewPurchaseDAO() *PurchaseDAO {
	return &PurchaseDAO{db: config.NewDatabaseConfig().DB}
}

// usado nos testes
func NewPurchaseDAOWithDB(db *sql.DB) *PurchaseDAO {
	return &PurchaseDAO{db: db}
}

func (d *PurchaseDAO) Create(p *model.Purchase) error {

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO purchase
		(invoice_number, invoice_type, supplier_id, status, total)
		VALUES (?, ?, ?, ?, ?)
	`

	res, err := tx.Exec(
		query,
		p.InvoiceNumber,
		p.InvoiceType,
		p.SupplierID,
		p.Status,
		p.Total,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	purchaseID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	p.ID = purchaseID

	for _, item := range p.Items {

		itemQuery := `
			INSERT INTO purchase_item
			(purchase_id, product_id, quantity, cost_price, total, cost_center_id)
			VALUES (?, ?, ?, ?, ?, ?)
		`

		_, err := tx.Exec(
			itemQuery,
			purchaseID,
			item.ProductID,
			item.Quantity,
			item.CostPrice,
			item.Total,
			item.CostCenterID,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (d *PurchaseDAO) FindByID(id int64) (*model.Purchase, error) {

	query := `
		SELECT id, invoice_number, invoice_type, supplier_id, status, total, created_at
		FROM purchase
		WHERE id = ?
	`

	var p model.Purchase

	err := d.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.InvoiceNumber,
		&p.InvoiceType,
		&p.SupplierID,
		&p.Status,
		&p.Total,
		&p.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	itemsQuery := `
		SELECT id, purchase_id, product_id, quantity, cost_price, total, cost_center_id
		FROM purchase_item
		WHERE purchase_id = ?
	`

	rows, err := d.db.Query(itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.PurchaseItem{}

	for rows.Next() {
		var item model.PurchaseItem

		err := rows.Scan(
			&item.ID,
			&item.PurchaseID,
			&item.ProductID,
			&item.Quantity,
			&item.CostPrice,
			&item.Total,
			&item.CostCenterID,
		)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	p.Items = items

	return &p, nil
}

func (d *PurchaseDAO) FindAll() ([]model.Purchase, error) {

	query := `
		SELECT id, invoice_number, invoice_type, supplier_id, status, total, created_at
		FROM purchase
		ORDER BY created_at DESC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []model.Purchase{}

	for rows.Next() {
		var p model.Purchase

		err := rows.Scan(
			&p.ID,
			&p.InvoiceNumber,
			&p.InvoiceType,
			&p.SupplierID,
			&p.Status,
			&p.Total,
			&p.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		list = append(list, p)
	}

	return list, nil
}

func (d *PurchaseDAO) Update(p *model.Purchase) error {

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	updateQuery := `
		UPDATE purchase
		SET invoice_number = ?, 
		    invoice_type = ?, 
		    supplier_id = ?, 
		    status = ?, 
		    total = ?
		WHERE id = ?
	`

	_, err = tx.Exec(
		updateQuery,
		p.InvoiceNumber,
		p.InvoiceType,
		p.SupplierID,
		p.Status,
		p.Total,
		p.ID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	// remove itens antigos
	deleteItemsQuery := `DELETE FROM purchase_item WHERE purchase_id = ?`

	_, err = tx.Exec(deleteItemsQuery, p.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// recria itens
	insertItemQuery := `
		INSERT INTO purchase_item
		(purchase_id, product_id, quantity, cost_price, total, cost_center_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	for _, item := range p.Items {

		_, err := tx.Exec(
			insertItemQuery,
			p.ID,
			item.ProductID,
			item.Quantity,
			item.CostPrice,
			item.Total,
			item.CostCenterID,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
