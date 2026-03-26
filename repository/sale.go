package repository

import (
	"database/sql"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
	"github.com/shopspring/decimal"
)

type PaymentMethodDAO interface {
	Create(p *model.PaymentMethod) error
	Update(p *model.PaymentMethod) error
	GetByID(id int64) (*model.PaymentMethod, error)
	List() ([]model.PaymentMethod, error)
	Delete(id int64) error
}

type paymentMethodDAO struct {
	db *sql.DB
}

func NewPaymentMethodDAO() PaymentMethodDAO {
	return &paymentMethodDAO{db: config.NewDatabaseConfig().DB}
}

func (d *paymentMethodDAO) Create(p *model.PaymentMethod) error {
	query := `INSERT INTO payment_method (name) VALUES (?)`

	result, err := d.db.Exec(query, p.Name)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	p.ID = id
	return nil
}

func (d *paymentMethodDAO) GetByID(id int64) (*model.PaymentMethod, error) {
	query := `SELECT id, name, active FROM payment_method WHERE id=?`

	var p model.PaymentMethod
	err := d.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Active)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (d *paymentMethodDAO) Update(p *model.PaymentMethod) error {
	query := `UPDATE payment_method SET name=?, active=? WHERE id=?`
	_, err := d.db.Exec(query, p.Name, p.Active, p.ID)
	return err
}

func (d *paymentMethodDAO) Delete(id int64) error {
	_, err := d.db.Exec(`DELETE FROM payment_method WHERE id=?`, id)
	return err
}

func (d *paymentMethodDAO) List() ([]model.PaymentMethod, error) {
	rows, err := d.db.Query(`SELECT id, name, active FROM payment_method`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PaymentMethod

	for rows.Next() {
		var p model.PaymentMethod
		rows.Scan(&p.ID, &p.Name, &p.Active)
		list = append(list, p)
	}

	return list, nil
}

type SalesOrderDAO interface {
	Create(s *model.SalesOrder) error
	Update(s *model.SalesOrder) error
	GetByID(id int64) (*model.SalesOrder, error)
	List() ([]model.SalesOrder, error)
	Delete(id int64) error
}

type salesOrderDAO struct {
	db *sql.DB
}

func NewSalesOrderDAO() SalesOrderDAO {
	return &salesOrderDAO{db: config.NewDatabaseConfig().DB}
}

func (d *salesOrderDAO) Create(s *model.SalesOrder) error {
	query := `
	INSERT INTO sales_order
	(price, discount, status, payment_method_id)
	VALUES (?, ?, ?, ?)
	`

	result, err := d.db.Exec(
		query,
		s.Price.StringFixed(2),
		s.Discount.StringFixed(2),
		s.Status,
		s.PaymentMethodID,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	s.ID = id
	return nil
}

func (d *salesOrderDAO) GetByID(id int64) (*model.SalesOrder, error) {
	query := `
	SELECT id, price, discount, status, payment_method_id, active, created_at
	FROM sales_order WHERE id=?
	`

	var s model.SalesOrder
	var price, discount string

	err := d.db.QueryRow(query, id).Scan(
		&s.ID,
		&price,
		&discount,
		&s.Status,
		&s.PaymentMethodID,
		&s.Active,
		&s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	s.Price, _ = decimal.NewFromString(price)
	s.Discount, _ = decimal.NewFromString(discount)

	return &s, nil
}

func (d *salesOrderDAO) Update(s *model.SalesOrder) error {
	query := `
	UPDATE sales_order
	SET price=?, discount=?, status=?, payment_method_id=?, active=?
	WHERE id=?
	`

	_, err := d.db.Exec(
		query,
		s.Price.String(),
		s.Discount.String(),
		s.Status,
		s.PaymentMethodID,
		s.Active,
		s.ID,
	)

	return err
}

func (d *salesOrderDAO) Delete(id int64) error {
	_, err := d.db.Exec(`DELETE FROM sales_order WHERE id=?`, id)
	return err
}

func (d *salesOrderDAO) List() ([]model.SalesOrder, error) {
	rows, err := d.db.Query(`
	SELECT id, price, discount, status, payment_method_id, active, created_at
	FROM sales_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.SalesOrder

	for rows.Next() {
		var s model.SalesOrder
		var price, discount string

		rows.Scan(
			&s.ID,
			&price,
			&discount,
			&s.Status,
			&s.PaymentMethodID,
			&s.Active,
			&s.CreatedAt,
		)

		s.Price, _ = decimal.NewFromString(price)
		s.Discount, _ = decimal.NewFromString(discount)

		list = append(list, s)
	}

	return list, nil
}

type SalesOrderItemDAO interface {
	Create(i *model.SalesOrderItem) error
	ListByOrder(orderID int64) ([]model.SalesOrderItem, error)
	DeleteByOrder(orderID int64) error
}

type salesOrderItemDAO struct {
	db *sql.DB
}

func NewSalesOrderItemDAO() SalesOrderItemDAO {
	return &salesOrderItemDAO{db: config.NewDatabaseConfig().DB}
}

func (d *salesOrderItemDAO) Create(i *model.SalesOrderItem) error {
	query := `
	INSERT INTO sales_order_item
	(sales_order_id, product_id, quantity, unit_price)
	VALUES (?, ?, ?, ?)
	`

	_, err := d.db.Exec(
		query,
		i.SalesOrderID,
		i.ProductID,
		i.Quantity,
		i.UnitPrice.StringFixed(2),
	)

	return err
}

func (d *salesOrderItemDAO) ListByOrder(orderID int64) ([]model.SalesOrderItem, error) {
	rows, err := d.db.Query(`
	SELECT sales_order_id, product_id, quantity, unit_price
	FROM sales_order_item
	WHERE sales_order_id=?`, orderID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.SalesOrderItem

	for rows.Next() {
		var i model.SalesOrderItem
		var price string

		rows.Scan(
			&i.SalesOrderID,
			&i.ProductID,
			&i.Quantity,
			&price,
		)

		i.UnitPrice, _ = decimal.NewFromString(price)
		list = append(list, i)
	}

	return list, nil
}

func (d *salesOrderItemDAO) DeleteByOrder(orderID int64) error {
	_, err := d.db.Exec(
		`DELETE FROM sales_order_item WHERE sales_order_id=?`,
		orderID,
	)
	return err
}
