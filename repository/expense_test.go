package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func newMockDAO(t *testing.T) (*expenseDAO, sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dao := &expenseDAO{db: db}

	return dao, mock, db
}

func TestExpenseDAO_Create(t *testing.T) {
	dao, mock, db := newMockDAO(t)
	defer db.Close()

	exp := &model.Expense{
		Description: "Internet",
		CategoryID:  1,
		Amount:      decimal.NewFromFloat(99.9),
		Date:        time.Now(),
		Status:      "paid",
		Attachments: []model.ExpenseAttachment{
			{URL: "url1"},
			{URL: "url2"},
		},
	}

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO expenses")).
		WillReturnResult(sqlmock.NewResult(10, 1))

	mock.ExpectExec("INSERT INTO expense_attachments").
		WithArgs(int64(10), "url1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO expense_attachments").
		WithArgs(int64(10), "url2").
		WillReturnResult(sqlmock.NewResult(2, 1))

	mock.ExpectCommit()

	err := dao.Create(exp)

	assert.NoError(t, err)
	assert.Equal(t, int64(10), exp.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExpenseDAO_Update(t *testing.T) {
	dao, mock, db := newMockDAO(t)
	defer db.Close()

	exp := &model.Expense{
		ID:          10,
		Description: "Updated",
		CategoryID:  2,
		Amount:      decimal.NewFromFloat(50),
		Date:        time.Now(),
		Status:      "pending",
		Attachments: []model.ExpenseAttachment{
			{URL: "new1"},
		},
	}

	mock.ExpectBegin()

	mock.ExpectExec("UPDATE expenses").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("DELETE FROM expense_attachments").
		WithArgs(exp.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("INSERT INTO expense_attachments").
		WithArgs(exp.ID, "new1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := dao.Update(exp)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExpenseDAO_Delete(t *testing.T) {
	dao, mock, db := newMockDAO(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM expenses").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Delete(1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExpenseDAO_AddAttachment(t *testing.T) {
	dao, mock, db := newMockDAO(t)
	defer db.Close()

	mock.ExpectExec("INSERT INTO expense_attachments").
		WithArgs(int64(1), "url").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := dao.AddAttachment(1, "url")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExpenseDAO_RemoveAttachment(t *testing.T) {
	dao, mock, db := newMockDAO(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM expense_attachments").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.RemoveAttachment(1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExpenseDAO_FindByID(t *testing.T) {
	dao, mock, db := newMockDAO(t)
	defer db.Close()

	now := time.Now()

	row := sqlmock.NewRows([]string{
		"id", "description", "category_id", "amount", "date", "status",
		"created_at", "updated_at",
		"id", "name", "created_at", "updated_at",
	}).AddRow(
		1, "Internet", 1, decimal.NewFromFloat(99),
		now, "paid", now, now,
		1, "Infra", now, now,
	)

	mock.ExpectQuery("SELECT").
		WithArgs(int64(1)).
		WillReturnRows(row)

	attRows := sqlmock.NewRows([]string{
		"id", "expense_id", "url", "created_at",
	}).AddRow(
		1, 1, "url", now,
	)

	mock.ExpectQuery("SELECT id, expense_id, url, created_at FROM expense_attachments").
		WithArgs(int64(1)).
		WillReturnRows(attRows)

	exp, err := dao.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, exp)
	assert.Len(t, exp.Attachments, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExpenseDAO_FindAllExpenseCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := &expenseDAO{db: db}

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"created_at",
		"updated_at",
	}).
		AddRow(1, "Alimentação", now, now).
		AddRow(2, "Higiene", now, now)

	mock.ExpectQuery(`SELECT id, name, created_at, updated_at FROM expenses_category ORDER BY name`).
		WillReturnRows(rows)

	result, err := dao.FindAllExpenseCategories()

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	assert.Equal(t, int64(1), result[0].ID)
	assert.Equal(t, "Alimentação", result[0].Name)

	assert.Equal(t, int64(2), result[1].ID)
	assert.Equal(t, "Higiene", result[1].Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}
