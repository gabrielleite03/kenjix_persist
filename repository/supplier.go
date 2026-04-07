package repository

import (
	"database/sql"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
)

type SupplierDAO interface {
	FindAll() ([]model.Supplier, error)
	FindByID(id int64) (*model.Supplier, error)
	Create(supplier *model.Supplier) error
	Update(supplier *model.Supplier) error
	Delete(id int64) error
}

type supplierDAO struct {
	db *sql.DB
}

func NewSupplierDAO(db *sql.DB) SupplierDAO {
	return &supplierDAO{db: config.NewDatabaseConfig().DB}
}

func (d *supplierDAO) FindAll() ([]model.Supplier, error) {
	query := `
		SELECT 
			id,
			razao_social,
			nome_fantasia,
			cnpj,
			ie,
			address,
			sales_person,
			email,
			phone,
			active,
			category_id
		FROM supplier
		ORDER BY razao_social
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []model.Supplier

	for rows.Next() {
		var s model.Supplier

		err := rows.Scan(
			&s.ID,
			&s.RazaoSocial,
			&s.NomeFantasia,
			&s.CNPJ,
			&s.IE,
			&s.Address,
			&s.Salesperson,
			&s.Email,
			&s.Phone,
			&s.Active,
			&s.CategoryID,
		)

		if err != nil {
			return nil, err
		}

		suppliers = append(suppliers, s)
	}

	return suppliers, nil
}

func (d *supplierDAO) FindByID(id int64) (*model.Supplier, error) {
	query := `
		SELECT 
			id,
			razao_social,
			nome_fantasia,
			cnpj,
			ie,
			address,
			sales_person,
			email,
			phone,
			active,
			category_id
		FROM supplier
		WHERE id = ?
	`

	var s model.Supplier

	err := d.db.QueryRow(query, id).Scan(
		&s.ID,
		&s.RazaoSocial,
		&s.NomeFantasia,
		&s.CNPJ,
		&s.IE,
		&s.Address,
		&s.Salesperson,
		&s.Email,
		&s.Phone,
		&s.Active,
		&s.CategoryID,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (d *supplierDAO) Create(s *model.Supplier) error {
	query := `
		INSERT INTO supplier (
			razao_social,
			nome_fantasia,
			cnpj,
			ie,
			address,
			sales_person,
			email,
			phone,
			active,
			category_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := d.db.Exec(
		query,
		s.RazaoSocial,
		s.NomeFantasia,
		s.CNPJ,
		s.IE,
		s.Address,
		s.Salesperson,
		s.Email,
		s.Phone,
		s.Active,
		s.CategoryID,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	s.ID = id
	return nil
}

func (d *supplierDAO) Update(s *model.Supplier) error {
	query := `
		UPDATE supplier SET
			razao_social = ?,
			nome_fantasia = ?,
			cnpj = ?,
			ie = ?,
			address = ?,
			sales_person = ?,
			email = ?,
			phone = ?,
			active = ?,
			category_id = ?
		WHERE id = ?
	`

	_, err := d.db.Exec(
		query,
		s.RazaoSocial,
		s.NomeFantasia,
		s.CNPJ,
		s.IE,
		s.Address,
		s.Salesperson,
		s.Email,
		s.Phone,
		s.Active,
		s.CategoryID,
		s.ID,
	)

	return err
}

func (d *supplierDAO) Delete(id int64) error {
	query := `DELETE FROM supplier WHERE id = ?`
	_, err := d.db.Exec(query, id)
	return err
}
