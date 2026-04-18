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

type StockMovementDAO struct {
	db *sql.DB
}

func NewStockMovementDAO() *StockMovementDAO {
	return &StockMovementDAO{
		db: config.NewDatabaseConfig().DB,
	}
}

func (d *StockDAO) Create(stock *model.Stock) error {

	query := `
	INSERT INTO stock (
		product_id,
		warehouse_place_id,
		purchase_item_id,
		quantity
	) VALUES (?, ?, ?, ?)
	`

	_, err := d.db.Exec(
		query,
		stock.Product.ID,
		stock.WarehousePlace.ID,
		stock.PurchaseItem.ID,
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
		purchase_item_id = ?,
		quantity = ?,
		active = ?
	WHERE id = ?
	`

	_, err := d.db.Exec(
		query,
		stock.Product.ID,
		stock.WarehousePlace.ID,
		stock.PurchaseItem.ID,
		stock.Quantity,
		stock.Active,
		stock.ID,
	)

	return err
}

func (d *StockDAO) GetByID(id int64) (*model.Stock, error) {

	query := `
	SELECT
		id,
		product_id,
		warehouse_place_id,
		purchase_item_id,
		quantity,
		active,
		updated_at
	FROM stock
	WHERE id = ?
	`

	row := d.db.QueryRow(query, id)

	var s model.Stock

	err := row.Scan(
		&s.ID,
		&s.Product.ID,
		&s.WarehousePlace.ID,
		&s.PurchaseItem.ID,
		&s.Quantity,
		&s.Active,
		&s.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (d *StockDAO) Get(
	productID,
	warehousePlaceID,
	purchaseItemID int64,
) (*model.Stock, error) {

	query := `
	SELECT
		id,
		quantity,
		active,
		updated_at
	FROM stock
	WHERE product_id = ?
	  AND warehouse_place_id = ?
	  AND purchase_item_id = ?
	`

	row := d.db.QueryRow(
		query,
		productID,
		warehousePlaceID,
		purchaseItemID,
	)

	var s model.Stock
	s.Product.ID = productID
	s.WarehousePlace.ID = warehousePlaceID
	s.PurchaseItem.ID = purchaseItemID

	err := row.Scan(
		&s.ID,
		&s.Quantity,
		&s.Active,
		&s.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (d *StockDAO) GetByProduct(productID int64) ([]model.Stock, error) {

	query := `
	SELECT
		id,
		product_id,
		warehouse_place_id,
		purchase_item_id,
		quantity,
		active,
		updated_at
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
			&s.Product.ID,
			&s.WarehousePlace.ID,
			&s.PurchaseItem.ID,
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

func (d *StockDAO) Deactivate(id int64) error {

	query := `
	UPDATE stock
	SET active = false
	WHERE id = ?
	`

	_, err := d.db.Exec(query, id)

	return err
}

func (d *StockDAO) GetAllActive() ([]model.Stock, error) {

	query := `
	SELECT
		id,
		product_id,
		warehouse_place_id,
		purchase_item_id,
		quantity,
		active,
		updated_at
	FROM stock
	WHERE active = true
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
			&s.Product.ID,
			&s.WarehousePlace.ID,
			&s.PurchaseItem.ID,
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

func (d *StockDAO) GetGroupedByProductAndWarehouse() ([]model.Stock, error) {

	query := `
	SELECT
		product_id,
		warehouse_place_id,
		SUM(quantity) as quantity
	FROM stock
	WHERE active = true
	GROUP BY product_id, warehouse_place_id 
	HAVING quantity > 0 
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
			&s.Product.ID,
			&s.WarehousePlace.ID,
			&s.Quantity,
		)
		if err != nil {
			return nil, err
		}

		s.Active = true

		list = append(list, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *StockDAO) GetGroupedByProduct() ([]model.Stock, error) {

	query := `
	SELECT
		product_id,
		SUM(quantity) as quantity
	FROM stock
	WHERE active = true
	GROUP BY product_id
	HAVING quantity > 0
	ORDER BY product_id
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
			&s.Product.ID,
			&s.Quantity,
		)
		if err != nil {
			return nil, err
		}

		s.Active = true

		list = append(list, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *StockMovementDAO) Create(m *model.StockMovement) error {

	query := `
	INSERT INTO stock_movement (
		product_id,
		warehouse_place_id,
		purchase_item_id,
		type,
		quantity,
		reference_id,
		reference_type,
		reason
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(
		query,
		m.ProductID,
		m.WarehousePlaceID,
		m.PurchseItemID,
		m.Type,
		m.Quantity,
		m.ReferenceID,
		m.ReferenceType,
		m.Reason,
	)

	return err
}

func (d *StockMovementDAO) GetByID(id int64) (*model.StockMovement, error) {

	query := `
	SELECT
		id,
		product_id,
		warehouse_place_id,
		purchase_item_id,
		type,
		quantity,
		reference_id,
		reference_type,
		reason,
		created_at
	FROM stock_movement
	WHERE id = ?
	`

	row := d.db.QueryRow(query, id)

	var m model.StockMovement

	err := row.Scan(
		&m.ID,
		&m.ProductID,
		&m.WarehousePlaceID,
		&m.PurchseItemID,
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

	return &m, nil
}

func (d *StockMovementDAO) GetByProduct(productID int64) ([]model.StockMovement, error) {

	query := `
	SELECT
		id,
		product_id,
		warehouse_place_id,
		purchase_item_id,
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
			&m.ID,
			&m.ProductID,
			&m.WarehousePlaceID,
			&m.PurchseItemID,
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

func (d *StockMovementDAO) GetByProductAndWarehouse(
	productID,
	warehousePlaceID int64,
) ([]model.StockMovement, error) {

	query := `
	SELECT
		id,
		product_id,
		warehouse_place_id,
		purchase_item_id,
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
			&m.ID,
			&m.ProductID,
			&m.WarehousePlaceID,
			&m.PurchseItemID,
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

func (d *StockMovementDAO) GetAll() ([]model.StockMovement, error) {

	query := `
	SELECT
		id,
		product_id,
		warehouse_place_id,
		purchase_item_id,
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
			&m.ID,
			&m.ProductID,
			&m.WarehousePlaceID,
			&m.PurchseItemID,
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

func (d *StockMovementDAO) GetByReference(referenceID int64) ([]model.StockMovement, error) {

	query := `
	SELECT
		id,
		product_id,
		warehouse_place_id,
		purchase_item_id,
		type,
		quantity,
		reference_id,
		reference_type,
		reason,
		created_at
	FROM stock_movement
	WHERE reference_id = ?
	ORDER BY created_at DESC
	`

	rows, err := d.db.Query(query, referenceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.StockMovement

	for rows.Next() {

		var m model.StockMovement

		err := rows.Scan(
			&m.ID,
			&m.ProductID,
			&m.WarehousePlaceID,
			&m.PurchseItemID,
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

func (d *StockMovementDAO) FindAllEager() ([]model.StockMovementEager, error) {

	query := `
	SELECT 
		sm.id,
		sm.type,
		sm.quantity,
		sm.reference_id,
		sm.reference_type,
		sm.reason,
		sm.created_at,

		-- PRODUCT
		p.id,
		p.name,
		p.sku,
		p.price,
		p.marca,
		p.description,
		p.active,
		p.volume,
		p.category_id,

		-- WAREHOUSE PLACE
		wp.id,
		wp.name,
		wp.active,
		wp.warehouse_place_type_id,
		wp.warehouse_id,
		wp.capacity,

		-- PURCHASE ITEM
		pi.id,
		pi.purchase_id,
		pi.product_id,
		pi.quantity,
		pi.cost_price,
		pi.total,
		pi.cost_center_id

	FROM stock_movement sm
	JOIN product p ON p.id = sm.product_id
	JOIN warehouse_place wp ON wp.id = sm.warehouse_place_id
	JOIN purchase_item pi ON pi.id = sm.purchase_item_id
	ORDER BY sm.created_at DESC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.StockMovementEager

	for rows.Next() {

		var sm model.StockMovementEager

		var (
			productCategoryID *int64
			wpTypeID          *int64
			wpWarehouseID     *int64
			wpCapacity        *int64
			piCostCenterID    *int64
		)

		err := rows.Scan(
			&sm.ID,
			&sm.Type,
			&sm.Quantity,
			&sm.ReferenceID,
			&sm.ReferenceType,
			&sm.Reason,
			&sm.CreatedAt,

			&sm.Product.ID,
			&sm.Product.Name,
			&sm.Product.SKU,
			&sm.Product.Price,
			&sm.Product.Marca,
			&sm.Product.Description,
			&sm.Product.Active,
			&sm.Product.Volume,
			&productCategoryID,

			&sm.WarehousePlace.ID,
			&sm.WarehousePlace.Name,
			&sm.WarehousePlace.Active,
			&wpTypeID,
			&wpWarehouseID,
			&wpCapacity,

			&sm.PurchaseItem.ID,
			&sm.PurchaseItem.PurchaseID,
			&sm.PurchaseItem.ProductID,
			&sm.PurchaseItem.Quantity,
			&sm.PurchaseItem.CostPrice,
			&sm.PurchaseItem.Total,
			&piCostCenterID,
		)

		if err != nil {
			return nil, err
		}

		sm.Product.CategoryID = productCategoryID
		sm.WarehousePlace.WarehousePlaceTypeID = wpTypeID
		sm.WarehousePlace.WarehouseID = wpWarehouseID
		sm.WarehousePlace.Capacity = wpCapacity
		sm.PurchaseItem.CostCenterID = piCostCenterID

		list = append(list, sm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
