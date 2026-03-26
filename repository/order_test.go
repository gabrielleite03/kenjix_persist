package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gabrielleite03/kenjix_domain/model"
)

func TestPurchaseStatusDAO_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &purchaseStatusDAO{db: db}

	status := &model.PurchaseStatus{
		Name: "OPEN",
	}

	mock.ExpectExec("INSERT INTO purchase_status").
		WithArgs(status.Name, status.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := dao.Create(status)

	if err != nil {
		t.Fatal(err)
	}

	if status.ID != 1 {
		t.Fatal("ID não atribuído")
	}
}

func TestPurchaseStatusDAO_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &purchaseStatusDAO{db: db}

	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "active",
	}).AddRow(1, "OPEN", nil, true)

	mock.ExpectQuery("SELECT (.+) FROM purchase_status").
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

func TestPurchaseOrderDAO_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &purchaseOrderDAO{db: db}

	order := &model.PurchaseOrder{
		PurchaseStatusID:   1,
		FiscalNumberTypeID: 2,
		SupplierID:         3,
		Status:             "OPEN",
	}

	mock.ExpectExec("INSERT INTO purchase_order").
		WithArgs(
			order.PurchaseStatusID,
			order.FiscalNumberTypeID,
			order.SupplierID,
			order.FiscalNumber,
			order.Status,
		).
		WillReturnResult(sqlmock.NewResult(10, 1))

	err := dao.Create(order)

	if err != nil {
		t.Fatal(err)
	}

	if order.ID != 10 {
		t.Fatal("id errado")
	}
}

func TestPurchaseOrderItemDAO_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &purchaseOrderItemDAO{db: db}

	item := &model.PurchaseOrderItem{
		PurchaseOrderID: 1,
		ProductID:       2,
		Quantity:        10,
	}

	mock.ExpectExec("INSERT INTO purchase_order_item").
		WithArgs(
			item.PurchaseOrderID,
			item.ProductID,
			item.Quantity,
			item.UnitPrice.String(),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := dao.Create(item)

	if err != nil {
		t.Fatal(err)
	}
}

func TestPurchaseOrderItemDAO_ListByOrder(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dao := &purchaseOrderItemDAO{db: db}

	rows := sqlmock.NewRows([]string{
		"purchase_order_id",
		"product_id",
		"quantity",
		"unit_price",
		"active",
	}).AddRow(1, 2, 5, "10.50", true)

	mock.ExpectQuery("SELECT (.+) FROM purchase_order_item").
		WithArgs(1).
		WillReturnRows(rows)

	list, err := dao.ListByOrder(1)

	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 1 {
		t.Fatal("lista vazia")
	}
}
