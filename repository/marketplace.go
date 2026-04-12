package repository

import (
	"context"
	"database/sql"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
)

type MarketplaceDAO interface {
	Create(ctx context.Context, m *model.Marketplace) error
	Update(ctx context.Context, m *model.Marketplace) error
	FindByID(ctx context.Context, id int64) (*model.Marketplace, error)
	FindAll(ctx context.Context) ([]model.Marketplace, error)
	Delete(ctx context.Context, id int64) error
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
