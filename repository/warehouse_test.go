package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*warehouseDAO, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro criando mock: %v", err)
	}

	dao := &warehouseDAO{db: db}

	cleanup := func() {
		db.Close()
	}

	return dao, mock, cleanup
}

func TestCreate(t *testing.T) {
	dao, mock, cleanup := setup(t)
	defer cleanup()

	w := &model.Warehouse{
		Name:    "Main",
		Address: "Rua A",
		Active:  true,
	}

	mock.ExpectQuery(`INSERT INTO warehouse`).
		WithArgs(w.Name, w.Address, w.Capacity, w.Active).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := dao.Create(w)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate(t *testing.T) {
	dao, mock, cleanup := setup(t)
	defer cleanup()

	w := &model.Warehouse{
		ID:      1,
		Name:    "Updated",
		Address: "Rua B",
		Active:  true,
	}

	mock.ExpectExec(`UPDATE warehouse`).
		WithArgs(w.Name, w.Address, w.Capacity, w.Active, w.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	result, err := dao.Update(w)

	assert.NoError(t, err)
	assert.Equal(t, w.Name, result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByID(t *testing.T) {
	dao, mock, cleanup := setup(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "name", "address", "capacity", "active",
	}).AddRow(1, "Main", "Rua A", nil, true)

	mock.ExpectQuery(`SELECT id, name, address, capacity, active FROM warehouse WHERE id=\$1`).
		WithArgs(1).
		WillReturnRows(rows)

	result, err := dao.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Main", result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAll(t *testing.T) {
	dao, mock, cleanup := setup(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "name", "address", "capacity", "active",
	}).
		AddRow(1, "Main", "Rua A", nil, true).
		AddRow(2, "Backup", "Rua B", nil, true)

	mock.ExpectQuery(`SELECT id, name, address, capacity, active FROM warehouse`).
		WillReturnRows(rows)

	result, err := dao.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Main", result[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete(t *testing.T) {
	dao, mock, cleanup := setup(t)
	defer cleanup()

	mock.ExpectExec(`DELETE FROM warehouse WHERE id=\$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Delete(1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
