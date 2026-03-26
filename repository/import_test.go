package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
)

func TestImportProcessDAO_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &importProcessDAO{db: db}

	now := time.Now()

	process := &model.ImportProcess{
		PurchaseOrderID: 1,
		Incoterm:        "FOB",
		ExchangeRate:    decimal.NewFromFloat(5.23),
		Status:          "OPEN",
		ArrivalDate:     &now,
	}

	mock.ExpectExec("INSERT INTO import_process").
		WithArgs(
			process.PurchaseOrderID,
			process.Incoterm,
			process.ExchangeRate.String(),
			process.Status,
			process.ArrivalDate,
		).
		WillReturnResult(sqlmock.NewResult(10, 1))

	err := dao.Create(process)

	if err != nil {
		t.Fatalf("erro create: %v", err)
	}

	if process.ID != 10 {
		t.Fatalf("id incorreto")
	}
}

func TestImportProcessDAO_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &importProcessDAO{db: db}

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id",
		"purchase_order_id",
		"incoterm",
		"exchange_rate",
		"status",
		"arrival_date",
	}).AddRow(1, 2, "FOB", "5.20", "OPEN", now)

	mock.ExpectQuery("SELECT (.+) FROM import_process").
		WithArgs(1).
		WillReturnRows(rows)

	result, err := dao.GetByID(1)

	if err != nil {
		t.Fatal(err)
	}

	if result.ID != 1 {
		t.Fatal("id errado")
	}
}

func TestImportCostDAO_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &importCostDAO{db: db}

	cost := &model.ImportCost{
		ImportProcessID: 1,
		Type:            "FREIGHT",
		Amount:          decimal.NewFromFloat(1000),
		Currency:        "USD",
	}

	mock.ExpectExec("INSERT INTO import_cost").
		WithArgs(
			cost.ImportProcessID,
			cost.Type,
			cost.Description,
			cost.Amount.String(),
			cost.Currency,
		).
		WillReturnResult(sqlmock.NewResult(5, 1))

	err := dao.Create(cost)
	if err != nil {
		t.Fatal(err)
	}
}

func TestImportCostAllocationDAO_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &importCostAllocationDAO{db: db}

	alloc := &model.ImportCostAllocation{
		ImportCostID:    1,
		ProductID:       2,
		AllocatedAmount: decimal.NewFromFloat(500),
	}

	mock.ExpectExec("INSERT INTO import_cost_allocation").
		WithArgs(
			alloc.ImportCostID,
			alloc.ProductID,
			alloc.AllocatedAmount.String(),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Create(alloc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestImportCostAllocationDAO_ListByCost(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &importCostAllocationDAO{db: db}

	rows := sqlmock.NewRows([]string{
		"import_cost_id",
		"product_id",
		"allocated_amount",
	}).AddRow(1, 2, "100.50")

	mock.ExpectQuery("SELECT (.+) FROM import_cost_allocation").
		WithArgs(1).
		WillReturnRows(rows)

	list, err := dao.ListByCost(1)

	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 1 {
		t.Fatal("lista vazia")
	}
}
