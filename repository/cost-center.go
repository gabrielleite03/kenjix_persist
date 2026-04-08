package repository

import (
	"database/sql"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
)

type CostCenterDAO struct {
	db *sql.DB
}

func NewCostCenterDAO() *CostCenterDAO {
	return &CostCenterDAO{db: config.NewDatabaseConfig().DB}
}

func (d *CostCenterDAO) Create(costCenter *model.CostCenter) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		INSERT INTO cost_center (
			name,
			code,
			description
		) VALUES (?, ?, ?)
	`

	result, err := tx.Exec(
		query,
		costCenter.Name,
		costCenter.Code,
		costCenter.Description,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	costCenter.ID = id

	if len(costCenter.Properties) > 0 {
		propQuery := `
			INSERT INTO cost_center_property (
				cost_center_id,
				name,
				value,
				type
			) VALUES (?, ?, ?, ?)
		`

		for _, p := range costCenter.Properties {
			_, err = tx.Exec(
				propQuery,
				costCenter.ID,
				p.Name,
				p.Value,
				p.Type,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (d *CostCenterDAO) Update(costCenter *model.CostCenter) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		UPDATE cost_center
		SET 
			name = ?,
			code = ?,
			description = ?,
			active = ?
		WHERE id = ?
	`

	_, err = tx.Exec(
		query,
		costCenter.Name,
		costCenter.Code,
		costCenter.Description,
		costCenter.Active,
		costCenter.ID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM cost_center_property WHERE cost_center_id = ?`,
		costCenter.ID,
	)
	if err != nil {
		return err
	}

	if len(costCenter.Properties) > 0 {
		propQuery := `
			INSERT INTO cost_center_property (
				cost_center_id,
				name,
				value,
				type
			) VALUES (?, ?, ?, ?)
		`

		for _, p := range costCenter.Properties {
			_, err = tx.Exec(
				propQuery,
				costCenter.ID,
				p.Name,
				p.Value,
				p.Type,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (d *CostCenterDAO) FindByID(id int64) (*model.CostCenter, error) {
	query := `
		SELECT id, name, code, description, active
		FROM cost_center
		WHERE id = ?
	`

	row := d.db.QueryRow(query, id)

	var costCenter model.CostCenter
	err := row.Scan(
		&costCenter.ID,
		&costCenter.Name,
		&costCenter.Code,
		&costCenter.Description,
		&costCenter.Active,
	)
	if err != nil {
		return nil, err
	}

	properties, err := d.getProperties(id)
	if err != nil {
		return nil, err
	}

	costCenter.Properties = properties

	return &costCenter, nil
}

func (d *CostCenterDAO) FindAll() ([]model.CostCenter, error) {
	query := `
		SELECT id, name, code, description, active
		FROM cost_center
		WHERE active = true
		ORDER BY name
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var costCenters []model.CostCenter

	for rows.Next() {
		var cc model.CostCenter

		err := rows.Scan(
			&cc.ID,
			&cc.Name,
			&cc.Code,
			&cc.Description,
			&cc.Active,
		)
		if err != nil {
			return nil, err
		}

		properties, err := d.getProperties(cc.ID)
		if err != nil {
			return nil, err
		}

		cc.Properties = properties
		costCenters = append(costCenters, cc)
	}

	return costCenters, nil
}

func (d *CostCenterDAO) Delete(id int64) error {
	query := `
		UPDATE cost_center
		SET active = false
		WHERE id = ?
	`

	_, err := d.db.Exec(query, id)
	return err
}

func (d *CostCenterDAO) getProperties(costCenterID int64) ([]model.CostCenterProperty, error) {
	query := `
		SELECT id, cost_center_id, name, value, type
		FROM cost_center_property
		WHERE cost_center_id = ?
	`

	rows, err := d.db.Query(query, costCenterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var properties []model.CostCenterProperty

	for rows.Next() {
		var p model.CostCenterProperty

		err := rows.Scan(
			&p.ID,
			&p.CostCenterID,
			&p.Name,
			&p.Value,
			&p.Type,
		)
		if err != nil {
			return nil, err
		}

		properties = append(properties, p)
	}

	return properties, nil
}
