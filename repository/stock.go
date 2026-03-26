package repository

import (
	"database/sql"
	"time"

	"github.com/gabrielleite03/kenjix_domain/model"
)

type StockDAO interface {
	Upsert(s *model.Stock) error
	Get(productID, warehouseID int64) (*model.Stock, error)
	List() ([]model.Stock, error)
	Delete(productID, warehouseID int64) error
}

type stockDAO struct {
	db *sql.DB
}

func NewStockDAO(db *sql.DB) StockDAO {
	return &stockDAO{db: db}
}

func (d *stockDAO) Upsert(s *model.Stock) error {
	query := `
	INSERT INTO stock
	(product_id, warehouse_id, quantity, active, updated_at)
	VALUES (?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	quantity = VALUES(quantity),
	active = VALUES(active),
	updated_at = VALUES(updated_at)
	`

	_, err := d.db.Exec(
		query,
		s.ProductID,
		s.WarehouseID,
		s.Quantity,
		s.Active,
		time.Now(),
	)

	return err
}

func (d *stockDAO) Get(productID, warehouseID int64) (*model.Stock, error) {
	query := `
	SELECT product_id, warehouse_id, quantity, active, updated_at
	FROM stock
	WHERE product_id=? AND warehouse_id=?
	`

	var s model.Stock

	err := d.db.QueryRow(query, productID, warehouseID).Scan(
		&s.ProductID,
		&s.WarehouseID,
		&s.Quantity,
		&s.Active,
		&s.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (d *stockDAO) List() ([]model.Stock, error) {
	rows, err := d.db.Query(`
	SELECT product_id, warehouse_id, quantity, active, updated_at
	FROM stock`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Stock

	for rows.Next() {
		var s model.Stock
		rows.Scan(
			&s.ProductID,
			&s.WarehouseID,
			&s.Quantity,
			&s.Active,
			&s.UpdatedAt,
		)

		list = append(list, s)
	}

	return list, nil
}

func (d *stockDAO) Delete(productID, warehouseID int64) error {
	_, err := d.db.Exec(
		`DELETE FROM stock WHERE product_id=? AND warehouse_id=?`,
		productID,
		warehouseID,
	)
	return err
}

type StockMovementDAO interface {
	Create(m *model.StockMovement) error
	ListByProduct(productID int64) ([]model.StockMovement, error)
	ListByWarehouse(warehouseID int64) ([]model.StockMovement, error)
}

type stockMovementDAO struct {
	db *sql.DB
}

func NewStockMovementDAO(db *sql.DB) StockMovementDAO {
	return &stockMovementDAO{db: db}
}

func (d *stockMovementDAO) Create(m *model.StockMovement) error {
	query := `
	INSERT INTO stock_movement
	(product_id, warehouse_id, type, quantity, created_at, reason)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := d.db.Exec(
		query,
		m.ProductID,
		m.WarehouseID,
		m.Type,
		m.Quantity,
		m.CreatedAt,
		m.Reason,
	)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	m.ID = id

	return nil
}

func (d *stockMovementDAO) ListByProduct(productID int64) ([]model.StockMovement, error) {
	rows, err := d.db.Query(`
	SELECT id, product_id, warehouse_id, type, quantity, created_at, reason
	FROM stock_movement
	WHERE product_id=?
	ORDER BY created_at DESC
	`, productID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.StockMovement

	for rows.Next() {
		var m model.StockMovement

		rows.Scan(
			&m.ID,
			&m.ProductID,
			&m.WarehouseID,
			&m.Type,
			&m.Quantity,
			&m.CreatedAt,
			&m.Reason,
		)

		list = append(list, m)
	}

	return list, nil
}

func (d *stockMovementDAO) ListByWarehouse(warehouseID int64) ([]model.StockMovement, error) {
	rows, err := d.db.Query(`
	SELECT id, product_id, warehouse_id, type, quantity, created_at, reason
	FROM stock_movement
	WHERE warehouse_id=?
	ORDER BY created_at DESC
	`, warehouseID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.StockMovement

	for rows.Next() {
		var m model.StockMovement

		rows.Scan(
			&m.ID,
			&m.ProductID,
			&m.WarehouseID,
			&m.Type,
			&m.Quantity,
			&m.CreatedAt,
			&m.Reason,
		)

		list = append(list, m)
	}

	return list, nil
}
