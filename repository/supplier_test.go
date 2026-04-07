package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, SupplierDAO) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}

	dao := &supplierDAO{db: db}
	return db, mock, dao
}

func TestSupplierDAO_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dao := &supplierDAO{db: db}

	rows := sqlmock.NewRows([]string{
		"id",
		"razao_social",
		"nome_fantasia",
		"cnpj",
		"ie",
		"address",
		"sales_person",
		"email",
		"phone",
		"active",
		"category_id",
		"id",
		"name",
		"description",
		"active",
	}).AddRow(
		1,
		"Fornecedor LTDA",
		"Fornecedor",
		"123456",
		nil,
		nil,
		nil,
		nil,
		nil,
		true,
		2,
		2,
		"Categoria Teste",
		"Descricao",
		true,
	)

	mock.ExpectQuery("SELECT (.+) FROM supplier").
		WillReturnRows(rows)

	result, err := dao.FindAll()

	require.NoError(t, err)
	require.Len(t, result, 1)

	s := result[0]

	assert.Equal(t, int64(1), s.ID)
	assert.Equal(t, "Fornecedor LTDA", s.RazaoSocial)
	assert.NotNil(t, s.Category)
	assert.Equal(t, int64(2), s.Category.ID)
	assert.Equal(t, "Categoria Teste", s.Category.Name)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSupplierDAO_FindByID(t *testing.T) {
	db, mock, dao := setupTest(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "razao_social", "nome_fantasia", "cnpj",
		"ie", "address", "sales_person", "email",
		"phone", "active", "category_id",
	}).AddRow(
		1, "Empresa LTDA", "Empresa", "123",
		nil, nil, nil, nil,
		nil, true, nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta("WHERE id = ?")).
		WithArgs(1).
		WillReturnRows(rows)

	result, err := dao.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
}

func TestSupplierDAO_Create(t *testing.T) {
	db, mock, dao := setupTest(t)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO supplier")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	supplier := &model.Supplier{
		RazaoSocial:  "Empresa",
		NomeFantasia: "Fantasia",
		CNPJ:         "123",
		Active:       true,
	}

	err := dao.Create(supplier)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), supplier.ID)
}

func TestSupplierDAO_Update(t *testing.T) {
	db, mock, dao := setupTest(t)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("UPDATE supplier SET")).
		WithArgs(
			"Empresa",
			"Fantasia",
			"123",
			nil,
			nil,
			nil,
			nil,
			nil,
			true,
			nil,
			int64(1),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	supplier := &model.Supplier{
		ID:           1,
		RazaoSocial:  "Empresa",
		NomeFantasia: "Fantasia",
		CNPJ:         "123",
		Active:       true,
	}

	err := dao.Update(supplier)

	assert.NoError(t, err)
}

func TestSupplierDAO_Delete(t *testing.T) {
	db, mock, dao := setupTest(t)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM supplier")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Delete(1)

	assert.NoError(t, err)
}
