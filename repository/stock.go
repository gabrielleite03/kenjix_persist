package repository

import (
	"database/sql"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
)

type StockDAO struct {
	db *sql.DB
}

func NewStockDAO() *StockDAO {
	return &StockDAO{db: config.NewDatabaseConfig().DB}
}

func (d *StockDAO) Create(stock *model.Stock) error {

	query := `
		INSERT INTO stock (
			product_id,
			warehouse_place_id,
			quantity
		) VALUES (?, ?, ?)
	`

	_, err := d.db.Exec(
		query,
		stock.ProductID,
		stock.WarehousePlaceID,
		stock.Quantity,
	)

	return err
}

func (d *StockDAO) Update(stock *model.Stock) error {

	query := `
		UPDATE stock
		SET
			product_id = ?,
			warehouse_place_id = ?,
			quantity = ?
		WHERE id = ?
	`

	_, err := d.db.Exec(
		query,
		stock.ProductID,
		stock.WarehousePlaceID,
		stock.Quantity,
		stock.ID,
	)

	return err
}

func (d *StockDAO) Get(productID, warehousePlaceID int64) (*model.Stock, error) {

	query := `
	SELECT id, product_id, warehouse_place_id, quantity, active, updated_at
	FROM stock
	WHERE product_id = ? AND warehouse_place_id = ?
	`

	row := d.db.QueryRow(query, productID, warehousePlaceID)

	var s model.Stock

	err := row.Scan(
		&s.ID,
		&s.ProductID,
		&s.WarehousePlaceID,
		&s.Quantity,
		&s.Active,
		&s.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (d *StockDAO) InsertMovement(m *model.StockMovement) error {

	query := `
	INSERT INTO stock_movement
	(product_id, warehouse_place_id, type, quantity, reference_id, reference_type, reason)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(
		query,
		m.ProductID,
		m.WarehousePlaceID,
		m.Type,
		m.Quantity,
		m.ReferenceID,
		m.ReferenceType,
		m.Reason,
	)

	return err
}

func (d *StockDAO) AddStock(m *model.StockMovement) error {

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := d.InsertMovementTx(tx, m); err != nil {
		return err
	}

	stock := &model.Stock{
		ProductID:        m.ProductID,
		WarehousePlaceID: m.WarehousePlaceID,
		Quantity:         m.Quantity,
		Active:           true,
	}

	if err := d.UpsertTx(tx, stock); err != nil {
		return err
	}

	return tx.Commit()
}

func (d *StockDAO) UpsertTx(tx *sql.Tx, stock *model.Stock) error {

	query := `
	INSERT INTO stock (product_id, warehouse_place_id, quantity, active)
	VALUES (?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	    quantity = quantity + VALUES(quantity),
	    active = VALUES(active)
	`

	_, err := tx.Exec(
		query,
		stock.ProductID,
		stock.WarehousePlaceID,
		stock.Quantity,
		stock.Active,
	)

	return err
}

func (d *StockDAO) InsertMovementTx(tx *sql.Tx, m *model.StockMovement) error {

	query := `
	INSERT INTO stock_movement
	(product_id, warehouse_place_id, type, quantity, reference_id, reference_type, reason)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := tx.Exec(
		query,
		m.ProductID,
		m.WarehousePlaceID,
		m.Type,
		m.Quantity,
		m.ReferenceID,
		m.ReferenceType,
		m.Reason,
	)

	return err
}

func (d *StockDAO) GetByProduct(productID int64) ([]model.Stock, error) {

	query := `
	SELECT id, product_id, warehouse_place_id, quantity, active, updated_at
	FROM stock
	WHERE product_id = ?
	`

	rows, err := d.db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Stock

	for rows.Next() {
		var s model.Stock

		err := rows.Scan(
			&s.ID,
			&s.ProductID,
			&s.WarehousePlaceID,
			&s.Quantity,
			&s.Active,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, s)
	}

	return list, nil
}

func (d *StockDAO) GetAll() ([]model.Stock, error) {

	query := `
	SELECT id, product_id, warehouse_place_id, quantity, active, updated_at
	FROM stock
	ORDER BY product_id, warehouse_place_id
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Stock

	for rows.Next() {
		var s model.Stock

		err := rows.Scan(
			&s.ID,
			&s.ProductID,
			&s.WarehousePlaceID,
			&s.Quantity,
			&s.Active,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *StockDAO) GetMovementsByProduct(productID int64) ([]model.StockMovement, error) {

	query := `
	SELECT 
	    product_id,
	    warehouse_place_id,
	    type,
	    quantity,
	    reference_id,
	    reference_type,
	    reason,
	    created_at
	FROM stock_movement
	WHERE product_id = ?
	ORDER BY created_at DESC
	`

	rows, err := d.db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.StockMovement

	for rows.Next() {
		var m model.StockMovement

		err := rows.Scan(
			&m.ProductID,
			&m.WarehousePlaceID,
			&m.Type,
			&m.Quantity,
			&m.ReferenceID,
			&m.ReferenceType,
			&m.Reason,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *StockDAO) GetAllMovements() ([]model.StockMovement, error) {

	query := `
	SELECT 
	    product_id,
	    warehouse_place_id,
	    type,
	    quantity,
	    reference_id,
	    reference_type,
	    reason,
	    created_at
	FROM stock_movement
	ORDER BY created_at DESC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.StockMovement

	for rows.Next() {
		var m model.StockMovement

		err := rows.Scan(
			&m.ProductID,
			&m.WarehousePlaceID,
			&m.Type,
			&m.Quantity,
			&m.ReferenceID,
			&m.ReferenceType,
			&m.Reason,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *StockDAO) GetMovementsByProductAndWarehouse(productID, warehousePlaceID int64) ([]model.StockMovement, error) {

	query := `
	SELECT 
	    product_id,
	    warehouse_place_id,
	    type,
	    quantity,
	    reference_id,
	    reference_type,
	    reason,
	    created_at
	FROM stock_movement
	WHERE product_id = ? 
	  AND warehouse_place_id = ?
	ORDER BY created_at DESC
	`

	rows, err := d.db.Query(query, productID, warehousePlaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.StockMovement

	for rows.Next() {
		var m model.StockMovement

		err := rows.Scan(
			&m.ProductID,
			&m.WarehousePlaceID,
			&m.Type,
			&m.Quantity,
			&m.ReferenceID,
			&m.ReferenceType,
			&m.Reason,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, m)
	}

	return list, nil
}
