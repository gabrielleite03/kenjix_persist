package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	model "github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
)

func setupExpenses(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *productDAO) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}

	dao := &productDAO{db: db}
	return db, mock, dao
}

func TestCreate(t *testing.T) {
	db, mock, dao := setupExpenses(t)
	defer db.Close()

	price := decimal.NewFromFloat(99.90)

	product := &model.Product{
		Name:        "Produto",
		SKU:         "SKU1",
		Price:       price,
		Marca:       "Marca",
		Description: "Descricao",
	}

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO product
		(name, sku, price, marca, description, category_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`)).
		WithArgs(
			product.Name,
			product.SKU,
			product.Price.String(),
			product.Marca,
			product.Description,
			product.CategoryID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := dao.Create(product)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if product.ID != 1 {
		t.Fatalf("id não foi setado")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectativas não atendidas: %v", err)
	}
}

func TestUpdateExpense(t *testing.T) {
	db, mock, dao := setupExpenses(t)
	defer db.Close()

	price := decimal.NewFromFloat(10)

	product := &model.Product{
		ID:          1,
		Name:        "Produto",
		SKU:         "SKU",
		Price:       price,
		Marca:       "Marca",
		Description: "Desc",
		Active:      true,
	}

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE product
		SET name=?, sku=?, price=?, marca=?, description=?, active=?, category_id=?
		WHERE id=?
	`)).
		WithArgs(
			product.Name,
			product.SKU,
			product.Price.String(),
			product.Marca,
			product.Description,
			product.Active,
			product.CategoryID,
			product.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM product_property").
		WithArgs(product.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM product_image").
		WithArgs(product.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM product_video").
		WithArgs(product.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err := dao.Update(product)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectativas não atendidas: %v", err)
	}
}

func TestGetByID(t *testing.T) {
	db, mock, dao := setupExpenses(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "name", "sku", "price", "marca", "description", "active", "category_id",
	}).AddRow(
		1,
		"Produto",
		"SKU",
		"15.50",
		"Marca",
		"Descricao",
		true,
		nil,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, sku, price, marca, description, active, category_id
		FROM product WHERE id=?
	`)).
		WithArgs(1).
		WillReturnRows(rows)

	mock.ExpectQuery("SELECT id, product_id, name, value").
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "name", "value"}))

	mock.ExpectQuery("SELECT id, product_id, url, position, is_primary").
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "url", "position", "is_primary"}))

	mock.ExpectQuery("SELECT id, product_id, url, provider").
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "url", "provider"}))

	p, err := dao.GetByID(1)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if p.ID != 1 {
		t.Fatalf("produto incorreto")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectativas não atendidas: %v", err)
	}
}

func TestDeleteExpense(t *testing.T) {
	db, mock, dao := setupExpenses(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM product WHERE id=").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Delete(1)
	if err != nil {
		t.Fatalf("erro ao deletar: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectativas não atendidas: %v", err)
	}
}

func TestList(t *testing.T) {
	db, mock, dao := setupExpenses(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "name", "sku", "price", "marca", "description", "active", "category_id",
	}).AddRow(
		1,
		"Produto",
		"SKU",
		"22.30",
		"Marca",
		"Descricao",
		true,
		nil,
	)

	mock.ExpectQuery("SELECT id, name, sku, price, marca, description, active, category_id").
		WillReturnRows(rows)

	mock.ExpectQuery("SELECT id, product_id, name, value").
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "name", "value"}))

	mock.ExpectQuery("SELECT id, product_id, url, position, is_primary").
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "url", "position", "is_primary"}))

	mock.ExpectQuery("SELECT id, product_id, url, provider").
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "url", "provider"}))

	list, err := dao.List()
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("esperado 1 produto")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectativas não atendidas: %v", err)
	}
}
