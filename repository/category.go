package repository

import (
	"database/sql"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
)

type CategoryDAO interface {
	Create(category *model.Category) error
	Update(category *model.Category) error
	GetByID(id int64) (*model.Category, error)
	Delete(id int64) error
	List() ([]model.Category, error)
}

type categoryDAO struct {
	db *sql.DB
}

func NewCategoryDAO() CategoryDAO {
	return &categoryDAO{db: config.NewDatabaseConfig().DB}
}

func (d *categoryDAO) Create(c *model.Category) error {
	query := `
		INSERT INTO category (name, description)
		VALUES (?, ?)
	`

	result, err := d.db.Exec(
		query,
		c.Name,
		c.Description,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	c.ID = id
	return nil
}

func (d *categoryDAO) Update(c *model.Category) error {
	query := `
		UPDATE category
		SET name=?, description=?, active=?
		WHERE id=?
	`

	_, err := d.db.Exec(
		query,
		c.Name,
		c.Description,
		c.Active,
		c.ID,
	)

	return err
}

func (d *categoryDAO) GetByID(id int64) (*model.Category, error) {
	query := `
		SELECT id, name, description, active
		FROM category
		WHERE id=?
	`

	var c model.Category

	err := d.db.QueryRow(query, id).Scan(
		&c.ID,
		&c.Name,
		&c.Description,
		&c.Active,
	)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (d *categoryDAO) Delete(id int64) error {
	query := `DELETE FROM category WHERE id=?`
	_, err := d.db.Exec(query, id)
	return err
}

func (d *categoryDAO) List() ([]model.Category, error) {
	query := `
		SELECT id, name, description, active
		FROM category where active=true
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Category

	for rows.Next() {
		var c model.Category

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&c.Active,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, c)
	}

	return list, nil
}
