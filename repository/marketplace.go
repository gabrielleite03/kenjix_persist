package repository

import (
	"context"
	"database/sql"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
	"github.com/shopspring/decimal"
)

type MarketplaceDAO interface {
	Create(ctx context.Context, m *model.Marketplace) error
	Update(ctx context.Context, m *model.Marketplace) error
	FindByID(ctx context.Context, id int64) (*model.Marketplace, error)
	FindAll(ctx context.Context) ([]model.Marketplace, error)
	Delete(ctx context.Context, id int64) error
	CreateProductMarketplace(ctx context.Context, pm *model.ProductMarketplace) error
	UpdateProductMarketplace(ctx context.Context, pm *model.ProductMarketplace) error
	FindProductMarketplaceByID(ctx context.Context, id int64) (*model.ProductMarketplace, error)
	FindAllProductMarketplace(ctx context.Context) ([]model.ProductMarketplace, error)
}

type marketplaceDAO struct {
	db *sql.DB
}

func NewMarketplaceDAO() MarketplaceDAO {
	return &marketplaceDAO{db: config.NewDatabaseConfig().DB}
}

func (r *marketplaceDAO) Create(ctx context.Context, m *model.Marketplace) error {
	query := `
		INSERT INTO marketplace (
			name, logo, status, commission_rate, integration_type,
			api_url, api_key, api_secret, api_endpoint,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		m.Name,
		m.Logo,
		m.Status,
		m.CommissionRate,
		m.IntegrationType,
		m.APIURL,
		m.APIKey,
		m.APISecret,
		m.APIEndpoint,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	m.ID = id
	return nil
}

func (r *marketplaceDAO) Update(ctx context.Context, m *model.Marketplace) error {
	query := `
		UPDATE marketplace SET
			name = ?,
			logo = ?,
			status = ?,
			commission_rate = ?,
			integration_type = ?,
			api_url = ?,
			api_key = ?,
			api_secret = ?,
			api_endpoint = ?,
			updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		m.Name,
		m.Logo,
		m.Status,
		m.CommissionRate,
		m.IntegrationType,
		m.APIURL,
		m.APIKey,
		m.APISecret,
		m.APIEndpoint,
		m.ID,
	)

	return err
}

func (r *marketplaceDAO) FindByID(ctx context.Context, id int64) (*model.Marketplace, error) {
	query := `
		SELECT 
			id, name, logo, status, commission_rate, integration_type,
			api_url, api_key, api_secret, api_endpoint,
			created_at, updated_at, deleted_at
		FROM marketplace
		WHERE id = ? AND deleted_at IS NULL
	`

	var m model.Marketplace

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID,
		&m.Name,
		&m.Logo,
		&m.Status,
		&m.CommissionRate,
		&m.IntegrationType,
		&m.APIURL,
		&m.APIKey,
		&m.APISecret,
		&m.APIEndpoint,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &m, nil
}

func (r *marketplaceDAO) FindAll(ctx context.Context) ([]model.Marketplace, error) {
	query := `
		SELECT 
			id, name, logo, status, commission_rate, integration_type,
			api_url, api_key, api_secret, api_endpoint,
			created_at, updated_at, deleted_at
		FROM marketplace
		WHERE deleted_at IS NULL
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Marketplace

	for rows.Next() {
		var m model.Marketplace

		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Logo,
			&m.Status,
			&m.CommissionRate,
			&m.IntegrationType,
			&m.APIURL,
			&m.APIKey,
			&m.APISecret,
			&m.APIEndpoint,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, m)
	}

	return list, nil
}

func (r *marketplaceDAO) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE marketplace
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *marketplaceDAO) CreateProductMarketplace(ctx context.Context, pm *model.ProductMarketplace) error {
	query := `
		INSERT INTO product_marketplace (
			product_id,
			marketplace_id,
			external_id,
			product_url,
			price,
			listing_type,
			status,
			active,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	var price interface{}
	if pm.Price != nil {
		price = pm.Price.String()
	}

	var externalID interface{}
	if pm.ExternalID != nil {
		externalID = *pm.ExternalID
	}

	var listingType interface{}
	if pm.ListingType != nil {
		listingType = *pm.ListingType
	}

	var status interface{}
	if pm.Status != nil {
		status = *pm.Status
	}

	result, err := r.db.ExecContext(
		ctx,
		query,
		pm.ProductID,
		pm.MarketplaceID,
		externalID,
		pm.ProductURL,
		price,
		listingType,
		status,
		pm.Active,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	pm.ID = id
	return nil
}

func (r *marketplaceDAO) UpdateProductMarketplace(ctx context.Context, pm *model.ProductMarketplace) error {
	query := `
		UPDATE product_marketplace SET
			product_id = ?,
			marketplace_id = ?,
			external_id = ?,
			product_url = ?,
			price = ?,
			listing_type = ?,
			status = ?,
			active = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	var price interface{}
	if pm.Price != nil {
		price = pm.Price.String()
	}

	var externalID interface{}
	if pm.ExternalID != nil {
		externalID = *pm.ExternalID
	}

	var listingType interface{}
	if pm.ListingType != nil {
		listingType = *pm.ListingType
	}

	var status interface{}
	if pm.Status != nil {
		status = *pm.Status
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		pm.ProductID,
		pm.MarketplaceID,
		externalID,
		pm.ProductURL,
		price,
		listingType,
		status,
		pm.Active,
		pm.ID,
	)

	return err
}

func (r *marketplaceDAO) FindProductMarketplaceByID(ctx context.Context, id int64) (*model.ProductMarketplace, error) {
	query := `
		SELECT 
			id,
			product_id,
			marketplace_id,
			external_id,
			product_url,
			price,
			listing_type,
			status,
			active,
			created_at,
			updated_at
		FROM product_marketplace
		WHERE id = ? AND active = 1
	`

	var pm model.ProductMarketplace

	var externalID sql.NullString
	var listingType sql.NullString
	var status sql.NullString
	var price sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&pm.ID,
		&pm.ProductID,
		&pm.MarketplaceID,
		&externalID,
		&pm.ProductURL,
		&price,
		&listingType,
		&status,
		&pm.Active,
		&pm.CreatedAt,
		&pm.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// tratamentos de NULL
	if externalID.Valid {
		pm.ExternalID = &externalID.String
	}

	if listingType.Valid {
		pm.ListingType = &listingType.String
	}

	if status.Valid {
		pm.Status = &status.String
	}

	if price.Valid {
		d, err := decimal.NewFromString(price.String)
		if err != nil {
			return nil, err
		}
		pm.Price = &d
	}

	return &pm, nil
}

func (r *marketplaceDAO) FindAllProductMarketplace(ctx context.Context) ([]model.ProductMarketplace, error) {
	query := `
		SELECT 
			id,
			product_id,
			marketplace_id,
			external_id,
			product_url,
			price,
			listing_type,
			status,
			active,
			created_at,
			updated_at
		FROM product_marketplace
		WHERE active = 1
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.ProductMarketplace

	for rows.Next() {
		var pm model.ProductMarketplace

		var externalID sql.NullString
		var listingType sql.NullString
		var status sql.NullString
		var price sql.NullString

		err := rows.Scan(
			&pm.ID,
			&pm.ProductID,
			&pm.MarketplaceID,
			&externalID,
			&pm.ProductURL,
			&price,
			&listingType,
			&status,
			&pm.Active,
			&pm.CreatedAt,
			&pm.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// NULL handling
		if externalID.Valid {
			pm.ExternalID = &externalID.String
		}

		if listingType.Valid {
			pm.ListingType = &listingType.String
		}

		if status.Valid {
			pm.Status = &status.String
		}

		if price.Valid {
			d, err := decimal.NewFromString(price.String)
			if err != nil {
				return nil, err
			}
			pm.Price = &d
		}

		list = append(list, pm)
	}

	return list, nil
}
