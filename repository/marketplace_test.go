package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, MarketplaceDAO) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dao := &marketplaceDAO{db: db} // 🔥 injeção manual

	return db, mock, dao
}

func TestCreateMarketplace(t *testing.T) {
	db, mock, dao := setup(t)
	defer db.Close()

	m := &model.Marketplace{
		Name:            "Mercado Livre",
		Status:          "active",
		CommissionRate:  decimal.NewFromFloat(0.16),
		IntegrationType: "api",
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO marketplace (
			name, logo, status, commission_rate, integration_type,
			api_url, api_key, api_secret, api_endpoint,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`)).
		WithArgs(
			m.Name,
			m.Logo,
			m.Status,
			m.CommissionRate,
			m.IntegrationType,
			m.APIURL,
			m.APIKey,
			m.APISecret,
			m.APIEndpoint,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := dao.Create(context.Background(), m)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), m.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateMarketplace(t *testing.T) {
	db, mock, dao := setup(t)
	defer db.Close()

	m := &model.Marketplace{
		ID:              1,
		Name:            "Shopee",
		Status:          "active",
		CommissionRate:  decimal.NewFromFloat(0.18),
		IntegrationType: "api",
	}

	mock.ExpectExec(regexp.QuoteMeta(`
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
	`)).
		WithArgs(
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
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Update(context.Background(), m)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByID(t *testing.T) {
	db, mock, dao := setup(t)
	defer db.Close()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "logo", "status", "commission_rate", "integration_type",
		"api_url", "api_key", "api_secret", "api_endpoint",
		"created_at", "updated_at", "deleted_at",
	}).AddRow(
		1,
		"Amazon",
		nil,
		"inactive",
		"0.15",
		"manual",
		nil, nil, nil, nil,
		now,
		now,
		nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT 
			id, name, logo, status, commission_rate, integration_type,
			api_url, api_key, api_secret, api_endpoint,
			created_at, updated_at, deleted_at
		FROM marketplace
		WHERE id = ? AND deleted_at IS NULL
	`)).
		WithArgs(1).
		WillReturnRows(rows)

	result, err := dao.FindByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Amazon", result.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAll(t *testing.T) {
	db, mock, dao := setup(t)
	defer db.Close()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "logo", "status", "commission_rate", "integration_type",
		"api_url", "api_key", "api_secret", "api_endpoint",
		"created_at", "updated_at", "deleted_at",
	}).
		AddRow(1, "ML", nil, "active", "0.16", "api", nil, nil, nil, nil, now, now, nil).
		AddRow(2, "Shopee", nil, "active", "0.18", "api", nil, nil, nil, nil, now, now, nil)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT 
			id, name, logo, status, commission_rate, integration_type,
			api_url, api_key, api_secret, api_endpoint,
			created_at, updated_at, deleted_at
		FROM marketplace
		WHERE deleted_at IS NULL
		ORDER BY id DESC
	`)).
		WillReturnRows(rows)

	result, err := dao.FindAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteMarketplace(t *testing.T) {
	db, mock, dao := setup(t)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE marketplace
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Delete(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
