package repository

import (
	"database/sql"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
	"github.com/shopspring/decimal"
)

type ImportProcessDAO interface {
	Create(p *model.ImportProcess) error
	Update(p *model.ImportProcess) error
	GetByID(id int64) (*model.ImportProcess, error)
	Delete(id int64) error
	List() ([]model.ImportProcess, error)
}

type importProcessDAO struct {
	db *sql.DB
}

func NewImportProcessDAO() ImportProcessDAO {
	return &importProcessDAO{db: config.NewDatabaseConfig().DB}
}

func (d *importProcessDAO) Create(p *model.ImportProcess) error {
	query := `
		INSERT INTO import_process
		(purchase_order_id, incoterm, exchange_rate, status, arrival_date)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := d.db.Exec(
		query,
		p.PurchaseOrderID,
		p.Incoterm,
		p.ExchangeRate.String(),
		p.Status,
		p.ArrivalDate,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	p.ID = id
	return nil
}

func (d *importProcessDAO) GetByID(id int64) (*model.ImportProcess, error) {
	query := `
		SELECT id, purchase_order_id, incoterm, exchange_rate, status, arrival_date
		FROM import_process WHERE id=?
	`

	var p model.ImportProcess
	var rate string

	err := d.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.PurchaseOrderID,
		&p.Incoterm,
		&rate,
		&p.Status,
		&p.ArrivalDate,
	)
	if err != nil {
		return nil, err
	}

	p.ExchangeRate, _ = decimal.NewFromString(rate)
	return &p, nil
}

func (d *importProcessDAO) Update(p *model.ImportProcess) error {
	query := `
		UPDATE import_process
		SET purchase_order_id=?, incoterm=?, exchange_rate=?, status=?, arrival_date=?
		WHERE id=?
	`

	_, err := d.db.Exec(
		query,
		p.PurchaseOrderID,
		p.Incoterm,
		p.ExchangeRate.String(),
		p.Status,
		p.ArrivalDate,
		p.ID,
	)

	return err
}

func (d *importProcessDAO) Delete(id int64) error {
	_, err := d.db.Exec("DELETE FROM import_process WHERE id=?", id)
	return err
}

func (d *importProcessDAO) List() ([]model.ImportProcess, error) {
	query := `
		SELECT id, purchase_order_id, incoterm, exchange_rate, status, arrival_date
		FROM import_process
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.ImportProcess

	for rows.Next() {
		var p model.ImportProcess

		err := rows.Scan(
			&p.ID,
			&p.PurchaseOrderID,
			&p.Incoterm,
			&p.ExchangeRate,
			&p.Status,
			&p.ArrivalDate,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, p)
	}

	return list, nil
}

type ImportCostDAO interface {
	Create(c *model.ImportCost) error
	ListByProcess(processID int64) ([]model.ImportCost, error)
	Delete(id int64) error
}

type importCostDAO struct {
	db *sql.DB
}

func NewImportCostDAO() ImportCostDAO {
	return &importCostDAO{db: config.NewDatabaseConfig().DB}
}

func (d *importCostDAO) Create(c *model.ImportCost) error {
	query := `
		INSERT INTO import_cost
		(import_process_id, type, description, amount, currency)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := d.db.Exec(
		query,
		c.ImportProcessID,
		c.Type,
		c.Description,
		c.Amount.String(),
		c.Currency,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	c.ID = id
	return nil
}

func (d *importCostDAO) ListByProcess(processID int64) ([]model.ImportCost, error) {
	rows, err := d.db.Query(`
		SELECT id, import_process_id, type, description, amount, currency
		FROM import_cost
		WHERE import_process_id=?`, processID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.ImportCost

	for rows.Next() {
		var c model.ImportCost
		var amount string

		rows.Scan(
			&c.ID,
			&c.ImportProcessID,
			&c.Type,
			&c.Description,
			&amount,
			&c.Currency,
		)

		c.Amount, _ = decimal.NewFromString(amount)
		list = append(list, c)
	}

	return list, nil
}

func (d *importCostDAO) Delete(id int64) error {
	query := `DELETE FROM import_cost WHERE id=?`
	_, err := d.db.Exec(query, id)
	return err
}

type ImportCostAllocationDAO interface {
	Create(a *model.ImportCostAllocation) error
	ListByCost(costID int64) ([]model.ImportCostAllocation, error)
	DeleteByCost(costID int64) error
}

type importCostAllocationDAO struct {
	db *sql.DB
}

func NewImportCostAllocationDAO() ImportCostAllocationDAO {
	return &importCostAllocationDAO{db: config.NewDatabaseConfig().DB}
}

func (d *importCostAllocationDAO) Create(a *model.ImportCostAllocation) error {
	query := `
		INSERT INTO import_cost_allocation
		(import_cost_id, product_id, allocated_amount)
		VALUES (?, ?, ?)
	`

	_, err := d.db.Exec(
		query,
		a.ImportCostID,
		a.ProductID,
		a.AllocatedAmount.String(),
	)

	return err
}

func (d *importCostAllocationDAO) ListByCost(costID int64) ([]model.ImportCostAllocation, error) {
	rows, err := d.db.Query(`
		SELECT import_cost_id, product_id, allocated_amount
		FROM import_cost_allocation
		WHERE import_cost_id=?`, costID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.ImportCostAllocation

	for rows.Next() {
		var a model.ImportCostAllocation
		var amount string

		rows.Scan(
			&a.ImportCostID,
			&a.ProductID,
			&amount,
		)

		a.AllocatedAmount, _ = decimal.NewFromString(amount)
		list = append(list, a)
	}

	return list, nil
}

func (d *importCostAllocationDAO) DeleteByCost(costID int64) error {
	query := `DELETE FROM import_cost_allocation WHERE import_cost_id=?`
	_, err := d.db.Exec(query, costID)
	return err
}
