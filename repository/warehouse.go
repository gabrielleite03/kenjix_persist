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

	CreateWarehousePlaceType(w *model.WarehousePlaceType) (*model.WarehousePlaceType, error)
	UpdateWarehousePlaceType(w *model.WarehousePlaceType) (*model.WarehousePlaceType, error)
	FindByIDWarehousePlaceType(id int64) (*model.WarehousePlaceType, error)
	FindAllWarehousePlaceType() ([]*model.WarehousePlaceType, error)
	DeleteWarehousePlaceType(id int64) error

	CreateWarehousePlace(w *model.WarehousePlace) (*model.WarehousePlace, error)
	UpdateWarehousePlace(w *model.WarehousePlace) (*model.WarehousePlace, error)
	FindByIDWarehousePlace(id int64) (*model.WarehousePlace, error)
	FindByWarehouseID(warehouseID int64) ([]*model.WarehousePlace, error)
	FindAllWarehousePlace() ([]*model.WarehousePlace, error)
	DeleteWarehousePlace(id int64) error
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

	query := `INSERT INTO warehouse (name, address, capacity) VALUES (?, ?, ?)`

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

func (d *warehouseDAO) CreateWarehousePlaceType(w *model.WarehousePlaceType) (*model.WarehousePlaceType, error) {
	if w == nil {
		return nil, errors.New("warehouse place type is nil")
	}

	query := `INSERT INTO warehouse_place_type (name, value) VALUES (?, ?, ?)`

	result, err := d.db.Exec(query, w.Name, w.Value)
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

func (d *warehouseDAO) UpdateWarehousePlaceType(w *model.WarehousePlaceType) (*model.WarehousePlaceType, error) {
	if w == nil {
		return nil, errors.New("warehouse place type is nil")
	}

	query := `UPDATE warehouse_place_type 
	          SET name=?, value=?, active=? 
	          WHERE id=?`

	_, err := d.db.Exec(query, w.Name, w.Value, w.Active, w.ID)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (d *warehouseDAO) FindByIDWarehousePlaceType(id int64) (*model.WarehousePlaceType, error) {
	var w model.WarehousePlaceType

	query := `SELECT id, name, value, active 
	          FROM warehouse_place_type 
	          WHERE id=?`

	err := d.db.QueryRow(query, id).
		Scan(&w.ID, &w.Name, &w.Value, &w.Active)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &w, nil
}

func (d *warehouseDAO) FindAllWarehousePlaceType() ([]*model.WarehousePlaceType, error) {
	rows, err := d.db.Query(`
		SELECT id, name, value, active 
		FROM warehouse_place_type
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.WarehousePlaceType

	for rows.Next() {
		var w model.WarehousePlaceType
		err := rows.Scan(&w.ID, &w.Name, &w.Value, &w.Active)
		if err != nil {
			return nil, err
		}

		list = append(list, &w)
	}

	return list, nil
}

func (d *warehouseDAO) DeleteWarehousePlaceType(id int64) error {
	_, err := d.db.Exec(
		`DELETE FROM warehouse_place_type WHERE id=?`,
		id,
	)
	return err
}

func (d *warehouseDAO) CreateWarehousePlace(w *model.WarehousePlace) (*model.WarehousePlace, error) {
	if w == nil {
		return nil, errors.New("warehouse place is nil")
	}

	query := `
		INSERT INTO warehouse_place 
		(name, warehouse_place_type_id, warehouse_id, capacity)
		VALUES (?, ?, ?, ?)
	`

	result, err := d.db.Exec(
		query,
		w.Name,
		w.WarehousePlaceTypeID,
		w.WarehouseID,
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

func (d *warehouseDAO) UpdateWarehousePlace(w *model.WarehousePlace) (*model.WarehousePlace, error) {
	if w == nil {
		return nil, errors.New("warehouse place is nil")
	}

	query := `
		UPDATE warehouse_place 
		SET name=?, active=?, warehouse_place_type_id=?, capacity=?
		WHERE id=?
	`

	_, err := d.db.Exec(
		query,
		w.Name,
		w.Active,
		w.WarehousePlaceTypeID,
		w.Capacity,
		w.ID,
	)

	if err != nil {
		return nil, err
	}

	return w, nil
}

func (d *warehouseDAO) FindByIDWarehousePlace(id int64) (*model.WarehousePlace, error) {
	var w model.WarehousePlace

	query := `
		SELECT id, name, active, warehouse_place_type_id, warehouse_id, capacity
		FROM warehouse_place
		WHERE id=?
	`

	err := d.db.QueryRow(query, id).
		Scan(
			&w.ID,
			&w.Name,
			&w.Active,
			&w.WarehousePlaceTypeID,
			&w.WarehouseID,
			&w.Capacity,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &w, nil
}

func (d *warehouseDAO) FindByWarehouseID(warehouseID int64) ([]*model.WarehousePlace, error) {
	rows, err := d.db.Query(`
		SELECT id, name, active, warehouse_place_type_id, warehouse_id, capacity
		FROM warehouse_place
		WHERE warehouse_id = ?
	`, warehouseID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.WarehousePlace

	for rows.Next() {
		var w model.WarehousePlace

		err := rows.Scan(
			&w.ID,
			&w.Name,
			&w.Active,
			&w.WarehousePlaceTypeID,
			&w.WarehouseID,
			&w.Capacity,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, &w)
	}

	return list, nil
}

func (d *warehouseDAO) FindAllWarehousePlace() ([]*model.WarehousePlace, error) {
	rows, err := d.db.Query(`
		SELECT id, name, active, warehouse_place_type_id, warehouse_id, capacity
		FROM warehouse_place
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.WarehousePlace

	for rows.Next() {
		var w model.WarehousePlace

		err := rows.Scan(
			&w.ID,
			&w.Name,
			&w.Active,
			&w.WarehousePlaceTypeID,
			&w.WarehouseID,
			&w.Capacity,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, &w)
	}

	return list, nil
}

func (d *warehouseDAO) DeleteWarehousePlace(id int64) error {
	_, err := d.db.Exec(
		`DELETE FROM warehouse_place WHERE id=?`,
		id,
	)
	return err
}
