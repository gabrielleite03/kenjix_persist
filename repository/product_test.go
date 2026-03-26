package repository

import (
	"regexp"
	"testing"

	model "github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestProductDAO_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro sqlmock: %v", err)
	}
	defer db.Close()

	dao := &productDAO{db: db}

	product := &model.Product{
		Name:        "Produto Teste",
		SKU:         "SKU-001",
		Price:       decimal.NewFromFloat(100),
		Marca:       "Marca",
		Description: "Desc",
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

	err = dao.Create(product)
	if err != nil {
		t.Fatalf("erro create: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectativas não atendidas: %v", err)
	}
}

func TestProductDAO_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &productDAO{db: db}

	rows := sqlmock.NewRows([]string{
		"id", "name", "sku", "price", "marca", "description", "active", "category_id",
	}).AddRow(
		1,
		"Produto",
		"SKU",
		"199.90",
		"Marca",
		"Desc",
		true,
		nil,
	)

	mock.ExpectQuery("SELECT (.+) FROM product WHERE id=").
		WithArgs(1).
		WillReturnRows(rows)

	mock.ExpectQuery("SELECT (.+) FROM product_property").
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "name", "value"}))

	mock.ExpectQuery("SELECT (.+) FROM product_image").
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "url", "position", "is_primary"}))

	mock.ExpectQuery("SELECT (.+) FROM product_video").
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "url", "provider"}))

	product, err := dao.GetByID(1)
	if err != nil {
		t.Fatalf("erro get: %v", err)
	}

	if product.ID != 1 {
		t.Fatalf("id incorreto")
	}

	if !product.Price.Equal(decimal.RequireFromString("199.90")) {
		t.Fatalf("price incorreto")
	}
}

func TestProductDAO_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &productDAO{db: db}

	mock.ExpectExec("DELETE FROM product").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Delete(1)
	if err != nil {
		t.Fatalf("erro delete: %v", err)
	}
}

func TestProductDAO_List(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &productDAO{db: db}

	rows := sqlmock.NewRows([]string{
		"id", "name", "sku", "price", "marca", "description", "active", "category_id",
	}).AddRow(
		1, "Produto", "SKU", "10.00", "Marca", "Desc", true, nil,
	)

	mock.ExpectQuery("SELECT (.+) FROM product").
		WillReturnRows(rows)

	list, err := dao.List()
	if err != nil {
		t.Fatalf("erro list: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("esperado 1 item")
	}
}
