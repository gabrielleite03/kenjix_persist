package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestExpenseDAO_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := &expenseDAO{db: db}

	rows := sqlmock.NewRows([]string{
		"id", "description", "category_id", "amount", "date", "status",
		"created_at", "updated_at",
		"id", "name", "created_at", "updated_at",
	}).AddRow(
		"1",
		"Supermercado",
		1,
		decimal.NewFromFloat(100.50),
		time.Now(),
		"paid",
		time.Now(),
		time.Now(),
		1,
		"Alimentação",
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery("SELECT (.+) FROM expenses").
		WillReturnRows(rows)

	mock.ExpectQuery("SELECT id, expense_id, url, created_at FROM expense_attachments").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "expense_id", "url", "created_at",
		}))

	expenses, err := dao.FindAll()

	assert.NoError(t, err)
	assert.Len(t, expenses, 1)
	assert.Equal(t, "Supermercado", expenses[0].Description)
}

func TestExpenseDAO_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &expenseDAO{db: db}

	exp := &model.Expense{
		Description: "Internet",
		CategoryID:  1,
		Amount:      decimal.NewFromFloat(99.9),
		Date:        time.Now(),
		Status:      "paid",
	}

	mock.ExpectExec("INSERT INTO expenses").
		WithArgs(
			exp.Description,
			exp.CategoryID,
			exp.Amount,
			exp.Date,
			exp.Status,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := dao.Create(exp)

	assert.NoError(t, err)
}

func TestExpenseDAO_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &expenseDAO{db: db}

	mock.ExpectExec("DELETE FROM expenses").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := dao.Delete(1)

	assert.NoError(t, err)
}
