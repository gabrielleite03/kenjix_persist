package repository

import (
	"database/sql"
	"errors"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
)

type WarehouseDAO struct {
	db *sql.DB
}

func NewWarehouseDAO() *WarehouseDAO {
	return &WarehouseDAO{db: config.NewDatabaseConfig().DB}
}

// Create insere um novo warehouse
func (d *WarehouseDAO) Create(w *model.Warehouse) (*model.Warehouse, error) {
	if w == nil {
		return nil, errors.New("warehouse is nil")
	}

	query := `INSERT INTO warehouse (name, address, capacity, active) 
			  VALUES ($1, $2, $3, $4) RETURNING id`

	err := d.db.QueryRow(query, w.Name, w.Address, w.Capacity, w.Active).Scan(&w.ID)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Update atualiza um warehouse existente
func (d *WarehouseDAO) Update(w *model.Warehouse) (*model.Warehouse, error) {
	if w == nil {
		return nil, errors.New("warehouse is nil")
	}

	query := `UPDATE warehouse SET name=$1, address=$2, capacity=$3, active=$4 WHERE id=$5`
	_, err := d.db.Exec(query, w.Name, w.Address, w.Capacity, w.Active, w.ID)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// FindByID retorna um warehouse pelo ID
func (d *WarehouseDAO) FindByID(id int64) (*model.Warehouse, error) {
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
func (d *WarehouseDAO) FindAll() ([]*model.Warehouse, error) {
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
func (d *WarehouseDAO) Delete(id int64) error {
	_, err := d.db.Exec(`DELETE FROM warehouse WHERE id=$1`, id)
	return err
}
