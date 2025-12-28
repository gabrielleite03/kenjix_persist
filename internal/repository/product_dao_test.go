package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kenjix.com/persist/internal/model"
)

func TestProductDAO_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dao := &productDAO{db: db}
	ctx := context.Background()

	categoryID := int64(1)

	prod := &model.Product{
		Name:       "Mouse",
		SKU:        "SKU-123",
		Price:      100.0,
		Active:     true,
		CategoryID: &categoryID,
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(10)

	mock.ExpectQuery(`INSERT INTO product`).
		WithArgs(prod.Name, prod.SKU, prod.Price, prod.Active, prod.CategoryID).
		WillReturnRows(rows)

	id, err := dao.Create(ctx, prod)

	require.NoError(t, err)
	assert.Equal(t, int64(10), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductDAO_Create_NilProduct(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	dao := &productDAO{db: db}

	id, err := dao.Create(context.Background(), nil)

	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
}

func TestProductDAO_Create_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dao := &productDAO{db: db}
	ctx := context.Background()

	categoryID := int64(1)

	prod := &model.Product{
		Name:       "Mouse",
		SKU:        "SKU-123",
		Price:      100.0,
		Active:     true,
		CategoryID: &categoryID,
	}

	mock.ExpectQuery(`INSERT INTO product`).
		WithArgs(
			prod.Name,
			prod.SKU,
			prod.Price,
			prod.Active,
			prod.CategoryID,
		).
		WillReturnError(errors.New("db error"))

	id, err := dao.Create(ctx, prod)

	require.Error(t, err)
	assert.Equal(t, int64(0), id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductDAO_GetByID_Found(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &productDAO{db: db}
	ctx := context.Background()

	rows := sqlmock.NewRows(
		[]string{"id", "name", "sku", "price", "active", "category_id"},
	).AddRow(1, "Keyboard", "SKU-1", 200.0, true, 2)

	mock.ExpectQuery(`SELECT id, name, sku, price, active, category_id FROM product`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	prod, err := dao.GetByID(ctx, 1)

	require.NoError(t, err)
	require.NotNil(t, prod)
	assert.Equal(t, "Keyboard", prod.Name)
	assert.NotNil(t, prod.CategoryID)
	assert.Equal(t, int64(2), *prod.CategoryID)
}

func TestProductDAO_GetByID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &productDAO{db: db}

	mock.ExpectQuery(`SELECT id, name, sku, price, active, category_id FROM product`).
		WithArgs(int64(99)).
		WillReturnError(sql.ErrNoRows)

	prod, err := dao.GetByID(context.Background(), 99)

	require.NoError(t, err)
	assert.Nil(t, prod)
}

func TestProductDAO_GetByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dao := &productDAO{db: db}
	ctx := context.Background()

	productID := int64(10)

	mock.ExpectQuery(`SELECT id, name, sku, price, active, category_id FROM product WHERE id = \$1`).
		WithArgs(productID).
		WillReturnError(errors.New("db failure"))

	prod, err := dao.GetByID(ctx, productID)

	require.Error(t, err)
	assert.Nil(t, prod)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductDAO_GetByID_CategoryID_Null(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dao := &productDAO{db: db}
	ctx := context.Background()

	rows := sqlmock.NewRows(
		[]string{"id", "name", "sku", "price", "active", "category_id"},
	).
		AddRow(1, "Mouse", "SKU-1", 99.90, true, nil)

	mock.ExpectQuery(`SELECT id, name, sku, price, active, category_id FROM product WHERE id = \$1`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	prod, err := dao.GetByID(ctx, 1)

	require.NoError(t, err)
	require.NotNil(t, prod)
	assert.Nil(t, prod.CategoryID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductDAO_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &productDAO{db: db}

	prod := &model.Product{
		ID:     1,
		Name:   "Updated",
		SKU:    "SKU-UP",
		Price:  150,
		Active: true,
	}

	mock.ExpectExec(`UPDATE product SET`).
		WithArgs(prod.Name, prod.SKU, prod.Price, prod.Active, prod.CategoryID, prod.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	rows, err := dao.Update(context.Background(), prod)

	require.NoError(t, err)
	assert.Equal(t, int64(1), rows)
}

func TestProductDAO_Update_ProductIsNil(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dao := &productDAO{db: db}

	rows, err := dao.Update(context.Background(), nil)

	assert.Equal(t, int64(0), rows)
	require.Error(t, err)
	assert.EqualError(t, err, "product is nil")
}

func TestProductDAO_Update_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dao := &productDAO{db: db}

	categoryID := int64(1)
	prod := &model.Product{
		ID:         10,
		Name:       "Mouse",
		SKU:        "SKU-123",
		Price:      99.90,
		Active:     true,
		CategoryID: &categoryID,
	}

	dbErr := errors.New("db error")

	mock.ExpectExec(`UPDATE product SET`).
		WithArgs(
			prod.Name,
			prod.SKU,
			prod.Price,
			prod.Active,
			prod.CategoryID,
			prod.ID,
		).
		WillReturnError(dbErr)

	rows, err := dao.Update(context.Background(), prod)

	assert.Equal(t, int64(0), rows)
	require.Error(t, err)
	assert.Equal(t, dbErr, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductDAO_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &productDAO{db: db}

	mock.ExpectExec(`DELETE FROM product WHERE id`).
		WithArgs(int64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	rows, err := dao.Delete(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, int64(1), rows)
}

func TestProductDAO_Delete_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dao := &productDAO{db: db}

	dbErr := errors.New("delete error")

	mock.ExpectExec(`DELETE FROM product WHERE id=\$1`).
		WithArgs(int64(10)).
		WillReturnError(dbErr)

	rows, err := dao.Delete(context.Background(), 10)

	assert.Equal(t, int64(0), rows)
	require.Error(t, err)
	assert.Equal(t, dbErr, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductDAO_List_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dao := &productDAO{db: db}

	dbErr := errors.New("query error")

	mock.ExpectQuery(`SELECT id, name, sku, price, active, category_id FROM product`).
		WillReturnError(dbErr)

	result, err := dao.List(context.Background())

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Equal(t, dbErr, err)
}

func TestProductDAO_List_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dao := &productDAO{db: db}

	rows := sqlmock.NewRows(
		[]string{"id", "name", "sku", "price", "active", "category_id"},
	).
		AddRow("INVALID", "Mouse", "SKU-1", 10.0, true, 1) // id inválido

	mock.ExpectQuery(`SELECT id, name, sku, price, active, category_id FROM product`).
		WillReturnRows(rows)

	result, err := dao.List(context.Background())

	assert.Nil(t, result)
	require.Error(t, err)
}

func TestProductDAO_List_RowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dao := &productDAO{db: db}

	rows := sqlmock.NewRows(
		[]string{"id", "name", "sku", "price", "active", "category_id"},
	).
		AddRow(1, "Mouse", "SKU-1", 10.0, true, nil).
		RowError(0, errors.New("row iteration error"))

	mock.ExpectQuery(`SELECT id, name, sku, price, active, category_id FROM product`).
		WillReturnRows(rows)

	result, err := dao.List(context.Background())

	assert.Nil(t, result)
	require.Error(t, err)
}

func TestProductDAO_List(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &productDAO{db: db}

	rows := sqlmock.NewRows(
		[]string{"id", "name", "sku", "price", "active", "category_id"},
	).
		AddRow(1, "Mouse", "SKU-1", 50.0, true, nil).
		AddRow(2, "Keyboard", "SKU-2", 150.0, true, 3)

	mock.ExpectQuery(`SELECT id, name, sku, price, active, category_id FROM product`).
		WillReturnRows(rows)

	list, err := dao.List(context.Background())

	require.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Nil(t, list[0].CategoryID)
	assert.NotNil(t, list[1].CategoryID)
}

func TestNewProductRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(db)

	require.NotNil(t, repo)

	// garante que é a implementação correta
	dao, ok := repo.(*productDAO)
	require.True(t, ok, "expected ProductRepository to be *productDAO")

	// garante que o db foi injetado corretamente
	assert.Equal(t, db, dao.db)
}
