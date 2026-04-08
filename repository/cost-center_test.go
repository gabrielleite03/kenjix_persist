package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
)

func setupCC(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *CostCenterDAO) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}

	dao := &CostCenterDAO{db: db}
	return db, mock, dao
}

func TestCreateCostCenter(t *testing.T) {
	db, mock, dao := setupCC(t)
	defer db.Close()

	cc := &model.CostCenter{
		Name:        "Marketing",
		Code:        "MKT",
		Description: "Centro marketing",
	}

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO cost_center (
			name,
			code,
			description
		) VALUES (?, ?, ?)
	`)).
		WithArgs(cc.Name, cc.Code, cc.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := dao.Create(cc)
	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}

	if cc.ID != 1 {
		t.Errorf("ID esperado 1, obtido %d", cc.ID)
	}
}

func TestCreateCostCenterWithProperties(t *testing.T) {
	db, mock, dao := setupCC(t)
	defer db.Close()

	cc := &model.CostCenter{
		Name:        "Logistica",
		Code:        "LOG",
		Description: "Centro logistica",
		Properties: []model.CostCenterProperty{
			{
				Name:  "Frete",
				Value: decimal.NewFromFloat(10.5),
				Type:  model.CostCenterPropertyTypeValue,
			},
		},
	}

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO cost_center (
			name,
			code,
			description
		) VALUES (?, ?, ?)
	`)).
		WithArgs(cc.Name, cc.Code, cc.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta(`
			INSERT INTO cost_center_property (
				cost_center_id,
				name,
				value,
				type
			) VALUES (?, ?, ?, ?)
		`)).
		WithArgs(1, "Frete", cc.Properties[0].Value, model.CostCenterPropertyTypeValue).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := dao.Create(cc)
	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}
}

func TestUpdateCostCenter(t *testing.T) {
	db, mock, dao := setupCC(t)
	defer db.Close()

	cc := &model.CostCenter{
		ID:          1,
		Name:        "Financeiro",
		Code:        "FIN",
		Description: "Centro financeiro",
		Active:      true,
	}

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE cost_center
		SET 
			name = ?,
			code = ?,
			description = ?,
			active = ?
		WHERE id = ?
	`)).
		WithArgs(cc.Name, cc.Code, cc.Description, cc.Active, cc.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(
		regexp.QuoteMeta(`DELETE FROM cost_center_property WHERE cost_center_id = ?`),
	).
		WithArgs(cc.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := dao.Update(cc)
	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}
}

func TestFindCostCenterByID(t *testing.T) {
	db, mock, dao := setupCC(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "name", "code", "description", "active",
	}).AddRow(1, "Marketing", "MKT", "Centro", true)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, code, description, active
		FROM cost_center
		WHERE id = ?
	`)).
		WithArgs(1).
		WillReturnRows(rows)

	propRows := sqlmock.NewRows([]string{
		"id", "cost_center_id", "name", "value", "type",
	}).AddRow(1, 1, "Budget", decimal.NewFromFloat(1000), "value")

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, cost_center_id, name, value, type
		FROM cost_center_property
		WHERE cost_center_id = ?
	`)).
		WithArgs(1).
		WillReturnRows(propRows)

	cc, err := dao.FindByID(1)
	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}

	if cc.ID != 1 {
		t.Errorf("ID incorreto")
	}
}

func TestDeleteCostCenter(t *testing.T) {
	db, mock, dao := setupCC(t)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE cost_center
		SET active = false
		WHERE id = ?
	`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := dao.Delete(1)
	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}
}
