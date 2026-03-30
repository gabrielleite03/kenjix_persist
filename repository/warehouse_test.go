package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*WarehouseDAO, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}

	dao := &WarehouseDAO{db: db}

	cleanup := func() {
		db.Close()
	}

	return dao, mock, cleanup
}

func TestCreateWarehouse(t *testing.T) {
	dao, mock, cleanup := setupMockDB(t)
	defer cleanup()

	w := &model.Warehouse{
		Name:    "Central",
		Address: "Rua 1",
		Active:  true,
	}

	// Espera a query de insert e retorna ID
	mock.ExpectQuery(`INSERT INTO warehouse`).
		WithArgs(w.Name, w.Address, w.Capacity, w.Active).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := dao.Create(w)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)

	// Verifica se todas as expectativas do mock foram cumpridas
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByIDWarehouse(t *testing.T) {
	dao, mock, cleanup := setupMockDB(t)
	defer cleanup()

	expected := &model.Warehouse{
		ID:      1,
		Name:    "Central",
		Address: "Rua 1",
		Active:  true,
	}

	mock.ExpectQuery(`SELECT id, name, address, capacity, active FROM warehouse WHERE id=\$1`).
		WithArgs(expected.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "address", "capacity", "active"}).
			AddRow(expected.ID, expected.Name, expected.Address, expected.Capacity, expected.Active),
		)

	result, err := dao.FindByID(expected.ID)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Name, result.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteWarehouse(t *testing.T) {
	dao, mock, cleanup := setupMockDB(t)
	defer cleanup()

	id := int64(1)
	mock.ExpectExec(`DELETE FROM warehouse WHERE id=\$1`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Delete(id)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
