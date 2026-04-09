package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/stretchr/testify/assert"
)

func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, WarehouseDAO) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}

	dao := &warehouseDAO{db: db}
	return db, mock, dao
}

// =====================
// CREATE
// =====================
func TestWarehouse_Create(t *testing.T) {
	db, mock, dao := newMockDB(t)
	defer db.Close()

	capacity := int64(100)
	w := &model.Warehouse{
		Name:     "WH1",
		Address:  "Rua A",
		Capacity: &capacity,
	}

	mock.ExpectExec("INSERT INTO warehouse").
		WithArgs(w.Name, w.Address, w.Capacity).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := dao.Create(w)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// =====================
// UPDATE
// =====================
func TestWarehouse_Update(t *testing.T) {
	db, mock, dao := newMockDB(t)
	defer db.Close()

	capacity := int64(200)
	w := &model.Warehouse{
		ID:       1,
		Name:     "WH2",
		Address:  "Rua B",
		Capacity: &capacity,
		Active:   true,
	}

	mock.ExpectExec("UPDATE warehouse").
		WithArgs(w.Name, w.Address, w.Capacity, w.Active, w.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	result, err := dao.Update(w)

	assert.NoError(t, err)
	assert.Equal(t, w.Name, result.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// =====================
// FIND BY ID
// =====================
func TestWarehouse_FindByID(t *testing.T) {
	db, mock, dao := newMockDB(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "name", "address", "capacity", "active",
	}).AddRow(1, "WH1", "Rua A", 100, true)

	mock.ExpectQuery("SELECT id, name, address").
		WithArgs(1).
		WillReturnRows(rows)

	result, err := dao.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// =====================
// FIND ALL
// =====================
func TestWarehouse_FindAll(t *testing.T) {
	db, mock, dao := newMockDB(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "name", "address", "capacity", "active",
	}).
		AddRow(1, "WH1", "Rua A", 100, true).
		AddRow(2, "WH2", "Rua B", 200, false)

	mock.ExpectQuery("SELECT id, name, address").
		WillReturnRows(rows)

	result, err := dao.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// =====================
// DELETE
// =====================
func TestWarehouse_Delete(t *testing.T) {
	db, mock, dao := newMockDB(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM warehouse").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Delete(1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
