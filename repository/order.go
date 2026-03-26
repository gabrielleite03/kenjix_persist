package repository

import (
	"database/sql"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
	"github.com/shopspring/decimal"
)

type PurchaseStatusDAO interface {
	Create(p *model.PurchaseStatus) error
	Update(p *model.PurchaseStatus) error
	GetByID(id int64) (*model.PurchaseStatus, error)
	List() ([]model.PurchaseStatus, error)
	Delete(id int64) error
}

type purchaseStatusDAO struct {
	db *sql.DB
}

func NewPurchaseStatusDAO() PurchaseStatusDAO {
	return &purchaseStatusDAO{db: config.NewDatabaseConfig().DB}
}

func (d *purchaseStatusDAO) Create(p *model.PurchaseStatus) error {
	query := `INSERT INTO purchase_status (name, description)
	          VALUES (?, ?)`

	result, err := d.db.Exec(query, p.Name, p.Description)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	p.ID = id
	return nil
}

func (d *purchaseStatusDAO) GetByID(id int64) (*model.PurchaseStatus, error) {
	query := `SELECT id, name, description, active
	          FROM purchase_status WHERE id=?`

	var p model.PurchaseStatus
	err := d.db.QueryRow(query, id).
		Scan(&p.ID, &p.Name, &p.Description, &p.Active)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (d *purchaseStatusDAO) Update(p *model.PurchaseStatus) error {
	query := `UPDATE purchase_status
	          SET name=?, description=?, active=?
	          WHERE id=?`

	_, err := d.db.Exec(query, p.Name, p.Description, p.Active, p.ID)
	return err
}

func (d *purchaseStatusDAO) Delete(id int64) error {
	_, err := d.db.Exec(`DELETE FROM purchase_status WHERE id=?`, id)
	return err
}

func (d *purchaseStatusDAO) List() ([]model.PurchaseStatus, error) {
	rows, err := d.db.Query(`SELECT id, name, description, active FROM purchase_status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PurchaseStatus

	for rows.Next() {
		var p model.PurchaseStatus
		rows.Scan(&p.ID, &p.Name, &p.Description, &p.Active)
		list = append(list, p)
	}

	return list, nil
}

type FiscalNumberTypeDAO interface {
	List() ([]model.FiscalNumberType, error)
}

type fiscalNumberTypeDAO struct {
	db *sql.DB
}

func NewFiscalNumberTypeDAO() FiscalNumberTypeDAO {
	return &fiscalNumberTypeDAO{db: config.NewDatabaseConfig().DB}
}

func (d *fiscalNumberTypeDAO) List() ([]model.FiscalNumberType, error) {
	rows, err := d.db.Query(`
		SELECT id, name, description, active
		FROM fiscal_number_type`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.FiscalNumberType

	for rows.Next() {
		var f model.FiscalNumberType
		rows.Scan(&f.ID, &f.Name, &f.Description, &f.Active)
		list = append(list, f)
	}

	return list, nil
}

type PurchaseOrderDAO interface {
	Create(p *model.PurchaseOrder) error
	GetByID(id int64) (*model.PurchaseOrder, error)
	Update(p *model.PurchaseOrder) error
	Delete(id int64) error
	List() ([]model.PurchaseOrder, error)
}

type purchaseOrderDAO struct {
	db *sql.DB
}

func NewPurchaseOrderDAO() PurchaseOrderDAO {
	return &purchaseOrderDAO{db: config.NewDatabaseConfig().DB}
}

func (d *purchaseOrderDAO) Create(p *model.PurchaseOrder) error {
	query := `
	INSERT INTO purchase_order
	(purchase_status_id, fiscal_number_type_id, supplier_id,
	 fiscal_number, status)
	VALUES (?, ?, ?, ?, ?)
	`

	result, err := d.db.Exec(
		query,
		p.PurchaseStatusID,
		p.FiscalNumberTypeID,
		p.SupplierID,
		p.FiscalNumber,
		p.Status,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	p.ID = id
	return nil
}

func (d *purchaseOrderDAO) GetByID(id int64) (*model.PurchaseOrder, error) {
	query := `
	SELECT id, purchase_status_id, fiscal_number_type_id,
	       supplier_id, fiscal_number, status, created_at, active
	FROM purchase_order WHERE id=?
	`

	var p model.PurchaseOrder

	err := d.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.PurchaseStatusID,
		&p.FiscalNumberTypeID,
		&p.SupplierID,
		&p.FiscalNumber,
		&p.Status,
		&p.CreatedAt,
		&p.Active,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (d *purchaseOrderDAO) Update(p *model.PurchaseOrder) error {
	query := `
	UPDATE purchase_order
	SET purchase_status_id=?, fiscal_number_type_id=?,
	    supplier_id=?, fiscal_number=?, status=?, active=?
	WHERE id=?
	`

	_, err := d.db.Exec(
		query,
		p.PurchaseStatusID,
		p.FiscalNumberTypeID,
		p.SupplierID,
		p.FiscalNumber,
		p.Status,
		p.Active,
		p.ID,
	)

	return err
}

func (d *purchaseOrderDAO) Delete(id int64) error {
	_, err := d.db.Exec(`DELETE FROM purchase_order WHERE id=?`, id)
	return err
}

func (d *purchaseOrderDAO) List() ([]model.PurchaseOrder, error) {
	rows, err := d.db.Query(`
	SELECT id, purchase_status_id, fiscal_number_type_id,
	       supplier_id, fiscal_number, status, created_at, active
	FROM purchase_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PurchaseOrder

	for rows.Next() {
		var p model.PurchaseOrder
		rows.Scan(
			&p.ID,
			&p.PurchaseStatusID,
			&p.FiscalNumberTypeID,
			&p.SupplierID,
			&p.FiscalNumber,
			&p.Status,
			&p.CreatedAt,
			&p.Active,
		)

		list = append(list, p)
	}

	return list, nil
}

type PurchaseOrderItemDAO interface {
	Create(item *model.PurchaseOrderItem) error
	ListByOrder(orderID int64) ([]model.PurchaseOrderItem, error)
	DeleteByOrder(orderID int64) error
}

type purchaseOrderItemDAO struct {
	db *sql.DB
}

func NewPurchaseOrderItemDAO() PurchaseOrderItemDAO {
	return &purchaseOrderItemDAO{db: config.NewDatabaseConfig().DB}
}

func (d *purchaseOrderItemDAO) Create(i *model.PurchaseOrderItem) error {
	query := `
	INSERT INTO purchase_order_item
	(purchase_order_id, product_id, quantity, unit_price)
	VALUES (?, ?, ?, ?)
	`

	_, err := d.db.Exec(
		query,
		i.PurchaseOrderID,
		i.ProductID,
		i.Quantity,
		i.UnitPrice.String(),
	)

	return err
}

func (d *purchaseOrderItemDAO) ListByOrder(orderID int64) ([]model.PurchaseOrderItem, error) {
	rows, err := d.db.Query(`
	SELECT purchase_order_id, product_id, quantity, unit_price, active
	FROM purchase_order_item
	WHERE purchase_order_id=?`, orderID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PurchaseOrderItem

	for rows.Next() {
		var i model.PurchaseOrderItem
		var price string

		rows.Scan(
			&i.PurchaseOrderID,
			&i.ProductID,
			&i.Quantity,
			&price,
			&i.Active,
		)

		i.UnitPrice, _ = decimal.NewFromString(price)
		list = append(list, i)
	}

	return list, nil
}

func (d *purchaseOrderItemDAO) DeleteByOrder(orderID int64) error {
	_, err := d.db.Exec(
		`DELETE FROM purchase_order_item WHERE purchase_order_id=?`,
		orderID,
	)
	return err
}
