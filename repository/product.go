package repository

import (
	"database/sql"

	model "github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
	"github.com/shopspring/decimal"
)

type ProductDAO interface {
	Create(product *model.Product) error
	Update(product *model.Product) error
	GetByID(id int64) (*model.Product, error)
	Delete(id int64) error
	List() ([]model.Product, error)
}

type productDAO struct {
	db *sql.DB
}

func NewProductDAO() ProductDAO {
	return &productDAO{db: config.NewDatabaseConfig().DB}
}

func (d *productDAO) Create(p *model.Product) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO product
		(name, sku, price, marca, description, category_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(
		query,
		p.Name,
		p.SKU,
		p.Price.String(),
		p.Marca,
		p.Description,
		p.CategoryID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, _ := result.LastInsertId()
	p.ID = id

	if err := d.insertProperties(tx, p); err != nil {
		tx.Rollback()
		return err
	}

	if err := d.insertImages(tx, p); err != nil {
		tx.Rollback()
		return err
	}

	if err := d.insertVideos(tx, p); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (d *productDAO) Update(p *model.Product) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	query := `
		UPDATE product
		SET name=?, sku=?, price=?, marca=?, description=?, active=?, category_id=?
		WHERE id=?
	`

	_, err = tx.Exec(
		query,
		p.Name,
		p.SKU,
		p.Price.String(),
		p.Marca,
		p.Description,
		p.Active,
		p.CategoryID,
		p.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete children
	tx.Exec("DELETE FROM product_property WHERE product_id=?", p.ID)
	tx.Exec("DELETE FROM product_image WHERE product_id=?", p.ID)
	tx.Exec("DELETE FROM product_video WHERE product_id=?", p.ID)

	// reinsert
	d.insertProperties(tx, p)
	d.insertImages(tx, p)
	d.insertVideos(tx, p)

	return tx.Commit()
}

func (d *productDAO) GetByID(id int64) (*model.Product, error) {
	query := `
		SELECT id, name, sku, price, marca, description, active, category_id
		FROM product WHERE id=?
	`

	var p model.Product
	var price string

	err := d.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.SKU,
		&price,
		&p.Marca,
		&p.Description,
		&p.Active,
		&p.CategoryID,
	)
	if err != nil {
		return nil, err
	}

	p.Price, _ = decimal.NewFromString(price)

	d.loadProperties(&p)
	d.loadImages(&p)
	d.loadVideos(&p)

	return &p, nil
}

func (d *productDAO) Delete(id int64) error {
	_, err := d.db.Exec("DELETE FROM product WHERE id=?", id)
	return err
}

func (d *productDAO) List() ([]model.Product, error) {
	rows, err := d.db.Query(`
		SELECT id, name, sku, price, marca, description, active, category_id
		FROM product
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Product

	for rows.Next() {
		var p model.Product
		var price string

		rows.Scan(
			&p.ID,
			&p.Name,
			&p.SKU,
			&price,
			&p.Marca,
			&p.Description,
			&p.Active,
			&p.CategoryID,
		)

		p.Price, _ = decimal.NewFromString(price)

		d.loadProperties(&p)
		d.loadImages(&p)
		d.loadVideos(&p)

		list = append(list, p)
	}

	return list, nil
}

func (d *productDAO) insertProperties(tx *sql.Tx, p *model.Product) error {
	for _, prop := range p.Properties {
		_, err := tx.Exec(
			`INSERT INTO product_property (product_id, name, value)
			 VALUES (?, ?, ?)`,
			p.ID,
			prop.Name,
			prop.Value,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *productDAO) loadProperties(p *model.Product) {
	rows, _ := d.db.Query(
		`SELECT id, product_id, name, value 
		 FROM product_property WHERE product_id=?`, p.ID)

	defer rows.Close()

	for rows.Next() {
		var prop model.ProductProperty
		rows.Scan(&prop.ID, &prop.ProductID, &prop.Name, &prop.Value)
		p.Properties = append(p.Properties, prop)
	}
}

func (d *productDAO) insertImages(tx *sql.Tx, p *model.Product) error {
	for _, img := range p.Images {
		_, err := tx.Exec(
			`INSERT INTO product_image 
			(product_id, url, position, is_primary)
			VALUES (?, ?, ?, ?)`,
			p.ID,
			img.URL,
			img.Position,
			img.IsPrimary,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *productDAO) loadImages(p *model.Product) {
	rows, _ := d.db.Query(
		`SELECT id, product_id, url, position, is_primary
		 FROM product_image WHERE product_id=?`, p.ID)

	defer rows.Close()

	for rows.Next() {
		var img model.ProductImage
		rows.Scan(&img.ID, &img.ProductID, &img.URL, &img.Position, &img.IsPrimary)
		p.Images = append(p.Images, img)
	}
}

func (d *productDAO) insertVideos(tx *sql.Tx, p *model.Product) error {
	for _, v := range p.Videos {
		_, err := tx.Exec(
			`INSERT INTO product_video (product_id, url, provider)
			 VALUES (?, ?, ?)`,
			p.ID,
			v.URL,
			v.Provider,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *productDAO) loadVideos(p *model.Product) {
	rows, _ := d.db.Query(
		`SELECT id, product_id, url, provider
		 FROM product_video WHERE product_id=?`, p.ID)

	defer rows.Close()

	for rows.Next() {
		var v model.ProductVideo
		rows.Scan(&v.ID, &v.ProductID, &v.URL, &v.Provider)
		p.Videos = append(p.Videos, v)
	}
}
