package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	model "github.com/gabrielleite03/kenjix_domain/model"
)

func TestCategoryDAO_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer db.Close()

	dao := &categoryDAO{db: db}

	category := &model.Category{
		Name:        "Eletrônicos",
		Description: "Produtos eletrônicos",
	}

	mock.ExpectExec("INSERT INTO category").
		WithArgs(category.Name, category.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.Create(category)
	if err != nil {
		t.Fatalf("erro ao criar category: %v", err)
	}

	if category.ID != 1 {
		t.Fatalf("esperado ID 1, obtido %d", category.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectativas não atendidas: %v", err)
	}
}

func TestCategoryDAO_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &categoryDAO{db: db}

	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "active",
	}).AddRow(1, "Eletrônicos", "Produtos", true)

	mock.ExpectQuery("SELECT (.+) FROM category WHERE id=").
		WithArgs(1).
		WillReturnRows(rows)

	category, err := dao.GetByID(1)
	if err != nil {
		t.Fatalf("erro getById: %v", err)
	}

	if category.ID != 1 {
		t.Fatalf("id incorreto")
	}
}

func TestCategoryDAO_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &categoryDAO{db: db}

	category := &model.Category{
		ID:          1,
		Name:        "Atualizado",
		Description: "Nova descrição",
		Active:      true,
	}

	mock.ExpectExec("UPDATE category").
		WithArgs(
			category.Name,
			category.Description,
			category.Active,
			category.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Update(category)
	if err != nil {
		t.Fatalf("erro update: %v", err)
	}
}

func TestCategoryDAO_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &categoryDAO{db: db}

	mock.ExpectExec("DELETE FROM category").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Delete(1)
	if err != nil {
		t.Fatalf("erro delete: %v", err)
	}
}

func TestCategoryDAO_List(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &categoryDAO{db: db}

	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "active",
	}).AddRow(1, "Cat1", "Desc1", true).
		AddRow(2, "Cat2", "Desc2", false)

	mock.ExpectQuery("SELECT (.+) FROM category").
		WillReturnRows(rows)

	list, err := dao.List()
	if err != nil {
		t.Fatalf("erro list: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("esperado 2 registros")
	}
}
