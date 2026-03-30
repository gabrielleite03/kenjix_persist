package repository

import (
	"database/sql"
	"errors"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
)

type WarehouseDAO interface {
	Create(w *model.Warehouse) (*model.Warehouse, error)
	Update(w *model.Warehouse) (*model.Warehouse, error)
	FindByID(id int64) (*model.Warehouse, error)
	FindAll() ([]*model.Warehouse, error)
	Delete(id int64) error
}

type warehouseDAO struct {
	db *sql.DB
}

func NewWarehouseDAO() WarehouseDAO {
	return &warehouseDAO{db: config.NewDatabaseConfig().DB}
}

// Create insere um novo warehouse
func (d *warehouseDAO) Create(w *model.Warehouse) (*model.Warehouse, error) {
	if w == nil {
		return nil, errors.New("warehouse is nil")
	}

	query := `INSERT INTO warehouse (name, address, capacity, active) 
			  VALUES (?, ?, ?)`

	result, err := d.db.Exec(
		query,
		w.Name,
		w.Address,
		w.Capacity,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	w.ID = id

	return w, nil
}

// Update atualiza um warehouse existente
func (d *warehouseDAO) Update(w *model.Warehouse) (*model.Warehouse, error) {
	if w == nil {
		return nil, errors.New("warehouse is nil")
	}

	query := `UPDATE warehouse SET name=?, address=?, capacity=?, active=? WHERE id=?`

	_, err := d.db.Exec(
		query,
		w.Name,
		w.Address,
		w.Capacity,
		w.Active,
		w.ID,
	)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// FindByID retorna um warehouse pelo ID
func (d *warehouseDAO) FindByID(id int64) (*model.Warehouse, error) {
	var w model.Warehouse
	query := `SELECT id, name, address, capacity, active FROM warehouse WHERE id=$1`
	err := d.db.QueryRow(query, id).Scan(&w.ID, &w.Name, &w.Address, &w.Capacity, &w.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &w, nil
}

// FindAll retorna todos os warehouses
func (d *warehouseDAO) FindAll() ([]*model.Warehouse, error) {
	rows, err := d.db.Query(`SELECT id, name, address, capacity, active FROM warehouse`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []*model.Warehouse
	for rows.Next() {
		var w model.Warehouse
		if err := rows.Scan(&w.ID, &w.Name, &w.Address, &w.Capacity, &w.Active); err != nil {
			return nil, err
		}
		warehouses = append(warehouses, &w)
	}

	return warehouses, nil
}

// Delete remove um warehouse pelo ID
func (d *warehouseDAO) Delete(id int64) error {
	_, err := d.db.Exec(`DELETE FROM warehouse WHERE id=$1`, id)
	return err
}
