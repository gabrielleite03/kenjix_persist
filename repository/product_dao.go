package repository

import (
	"context"
	"database/sql"
	"errors"

	model "github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
)

type ProductRepository interface {
	Create(ctx context.Context, prod *model.Product) (int64, error)
	CreateProductProperties(ctx context.Context, productId int64, props map[string]string) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Product, error)
	Update(ctx context.Context, prod *model.Product) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
	List(ctx context.Context) ([]*model.Product, error)
}

// productDAO provides CRUD operations for products
type productDAO struct {
	dbConnection *config.DatabaseConnection
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository() ProductRepository {
	return &productDAO{dbConnection: config.NewDatabaseConfig()}
}

// Create inserts a new product and returns the inserted ID
func (p *productDAO) Create(ctx context.Context, prod *model.Product) (int64, error) {
	if prod == nil {
		return 0, errors.New("product is nil")
	}
	var id int64
	err := p.dbConnection.DB.QueryRowContext(ctx,
		`INSERT INTO product (name, sku, price, active, category_id) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		prod.Name, prod.SKU, prod.Price, prod.Active, prod.CategoryID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetByID returns a product by id
func (p *productDAO) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	row := p.dbConnection.DB.QueryRowContext(ctx, `SELECT id, name, sku, price, active, category_id FROM product WHERE id = $1`, id)
	prod := &model.Product{}
	var cat sql.NullInt64
	if err := row.Scan(&prod.ID, &prod.Name, &prod.SKU, &prod.Price, &prod.Active, &cat); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if cat.Valid {
		prod.CategoryID = &cat.Int64
	} else {
		prod.CategoryID = nil
	}
	return prod, nil
}

// Update modifies an existing product. Returns rows affected.
func (p *productDAO) Update(ctx context.Context, prod *model.Product) (int64, error) {
	if prod == nil {
		return 0, errors.New("product is nil")
	}
	res, err := p.dbConnection.DB.ExecContext(ctx, `UPDATE product SET name=$1, sku=$2, price=$3, active=$4, category_id=$5 WHERE id=$6`,
		prod.Name, prod.SKU, prod.Price, prod.Active, prod.CategoryID, prod.ID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Delete removes a product by id. Returns rows affected.
func (p *productDAO) Delete(ctx context.Context, id int64) (int64, error) {
	res, err := p.dbConnection.DB.ExecContext(ctx, `DELETE FROM product WHERE id=$1`, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// List returns all products (simple implementation, no pagination)
func (p *productDAO) List(ctx context.Context) ([]*model.Product, error) {
	rows, err := p.dbConnection.DB.QueryContext(ctx, `SELECT id, name, sku, price, active, category_id FROM product`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*model.Product
	for rows.Next() {
		prod := &model.Product{}
		var cat sql.NullInt64
		if err := rows.Scan(&prod.ID, &prod.Name, &prod.SKU, &prod.Price, &prod.Active, &cat); err != nil {
			return nil, err
		}
		if cat.Valid {
			prod.CategoryID = &cat.Int64
		} else {
			prod.CategoryID = nil
		}
		out = append(out, prod)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateProductProperties implements [ProductRepository].
func (p *productDAO) CreateProductProperties(ctx context.Context, productId int64, props map[string]string) (int64, error) {
	for key, value := range props {
		_, err := insertProperty(p.dbConnection, productId, key, value)
		if err != nil {
			return 0, err
		}
	}
	return 0, nil
}

func insertProperty(dbConnection *config.DatabaseConnection, productId int64, key string, value string) (int64, error) {
	var id int64
	err := dbConnection.DB.QueryRow(`INSERT INTO product_property (product_id, key, value) VALUES ($1, $2, $3) RETURNING id`, productId, key, value).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
