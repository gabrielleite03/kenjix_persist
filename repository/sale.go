package repository

import (
	"database/sql"

	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/internal/config"
	"github.com/shopspring/decimal"
)

type PaymentMethodDAO interface {
	Create(p *model.PaymentMethod) error
	Update(p *model.PaymentMethod) error
	GetByID(id int64) (*model.PaymentMethod, error)
	List() ([]model.PaymentMethod, error)
	Delete(id int64) error
}

type paymentMethodDAO struct {
	db *sql.DB
}

func NewPaymentMethodDAO() PaymentMethodDAO {
	return &paymentMethodDAO{db: config.NewDatabaseConfig().DB}
}

func (d *paymentMethodDAO) Create(p *model.PaymentMethod) error {
	query := `INSERT INTO payment_method (name) VALUES (?)`

	result, err := d.db.Exec(query, p.Name)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	p.ID = id
	return nil
}

func (d *paymentMethodDAO) GetByID(id int64) (*model.PaymentMethod, error) {
	query := `SELECT id, name, active FROM payment_method WHERE id=?`

	var p model.PaymentMethod
	err := d.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Active)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (d *paymentMethodDAO) Update(p *model.PaymentMethod) error {
	query := `UPDATE payment_method SET name=?, active=? WHERE id=?`
	_, err := d.db.Exec(query, p.Name, p.Active, p.ID)
	return err
}

func (d *paymentMethodDAO) Delete(id int64) error {
	_, err := d.db.Exec(`DELETE FROM payment_method WHERE id=?`, id)
	return err
}

func (d *paymentMethodDAO) List() ([]model.PaymentMethod, error) {
	rows, err := d.db.Query(`SELECT id, name, active FROM payment_method`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PaymentMethod

	for rows.Next() {
		var p model.PaymentMethod
		rows.Scan(&p.ID, &p.Name, &p.Active)
		list = append(list, p)
	}

	return list, nil
}

type SalesOrderDAO interface {
	Create(s *model.SalesOrder) (*model.SalesOrder, error)
	Update(s *model.SalesOrder) error
	GetByID(id int64) (*model.SalesOrder, error)
	List() ([]model.SalesOrder, error)
	Delete(id int64) error
}

type salesOrderDAO struct {
	db *sql.DB
}

func NewSalesOrderDAO() SalesOrderDAO {
	return &salesOrderDAO{db: config.NewDatabaseConfig().DB}
}

func (d *salesOrderDAO) Create(s *model.SalesOrder) (*model.SalesOrder, error) {

	query := `
	INSERT INTO sales_order (
		price,
		discount,
		status,
		payment_method_id,
		active,
		customer_name,
		customer_document,
		marketplace_id,
		external_order_id,
		external_pack_id,
		payment_status,
		delivery_status
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := d.db.Exec(
		query,
		s.Price.StringFixed(2),
		s.Discount.StringFixed(2),
		s.Status,
		s.PaymentMethodID,
		s.Active,

		s.CustomerName,
		s.CustomerDocument,

		s.MarketplaceID,
		s.ExternalOrderID,
		s.ExternalPackID,

		s.PaymentStatus,
		s.DeliveryStatus,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	s.ID = id

	return s, nil
}

func (d *salesOrderDAO) GetByID(id int64) (*model.SalesOrder, error) {
	query := `
	SELECT
		so.id,
		so.price,
		so.discount,
		so.status,
		so.payment_method_id,
		so.active,
		so.created_at,
		so.customer_name,
		so.customer_document,
		so.marketplace_id,
		so.external_order_id,
		so.external_pack_id,
		so.payment_status,
		so.delivery_status,
		so.updated_at,

		pm.id,
		pm.name,
		pm.active,

		m.id,
		m.name,
		m.description,
		m.commission_rate,
		m.active
	FROM sales_order so
	INNER JOIN payment_method pm ON pm.id = so.payment_method_id
	LEFT JOIN marketplace m ON m.id = so.marketplace_id
	WHERE so.id = ?
	LIMIT 1
	`

	var s model.SalesOrder
	var price, discount string

	var customerName sql.NullString
	var customerDocument sql.NullString
	var marketplaceID sql.NullInt64
	var externalOrderID sql.NullString
	var externalPackID sql.NullString
	var paymentStatus sql.NullString
	var deliveryStatus sql.NullString
	var updatedAt sql.NullTime

	var paymentMethod model.PaymentMethod

	var marketplace model.Marketplace
	var marketplaceIDJoined sql.NullInt64
	var marketplaceName sql.NullString
	var marketplaceDescription sql.NullString
	var marketplaceCommissionRate string
	var marketplaceActive sql.NullBool

	err := d.db.QueryRow(query, id).Scan(
		&s.ID,
		&price,
		&discount,
		&s.Status,
		&s.PaymentMethodID,
		&s.Active,
		&s.CreatedAt,
		&customerName,
		&customerDocument,
		&marketplaceID,
		&externalOrderID,
		&externalPackID,
		&paymentStatus,
		&deliveryStatus,
		&updatedAt,

		&paymentMethod.ID,
		&paymentMethod.Name,
		&paymentMethod.Active,

		&marketplaceIDJoined,
		&marketplaceName,
		&marketplaceDescription,
		&marketplaceCommissionRate,
		&marketplaceActive,
	)
	if err != nil {
		return nil, err
	}

	s.Price, _ = decimal.NewFromString(price)
	s.Discount, _ = decimal.NewFromString(discount)

	if customerName.Valid {
		s.CustomerName = &customerName.String
	}
	if customerDocument.Valid {
		s.CustomerDocument = &customerDocument.String
	}
	if marketplaceID.Valid {
		s.MarketplaceID = &marketplaceID.Int64
	}
	if externalOrderID.Valid {
		s.ExternalOrderID = &externalOrderID.String
	}
	if externalPackID.Valid {
		s.ExternalPackID = &externalPackID.String
	}
	if paymentStatus.Valid {
		s.PaymentStatus = &paymentStatus.String
	}
	if deliveryStatus.Valid {
		s.DeliveryStatus = &deliveryStatus.String
	}
	if updatedAt.Valid {
		s.UpdatedAt = &updatedAt.Time
	}

	s.PaymentMethod = &paymentMethod

	if marketplaceIDJoined.Valid {
		marketplace.ID = marketplaceIDJoined.Int64
		marketplace.Name = marketplaceName.String

		commissionRate, err := decimal.NewFromString(marketplaceCommissionRate)
		if err == nil {
			marketplace.CommissionRate = commissionRate
		}

		s.Marketplace = &marketplace
	}

	invoice, err := d.GetInvoiceBySalesOrderID(s.ID)
	if err == nil {
		s.Invoice = invoice
	}

	return &s, nil
}
func (d *salesOrderDAO) GetInvoiceBySalesOrderID(salesOrderID int64) (*model.SalesOrderInvoice, error) {
	query := `
	SELECT
		id,
		sales_order_id,
		nfe_key,
		number,
		series,
		status,
		status_code,
		status_reason,
		protocol,
		issued_at,
		authorized_at,
		cancelled_at,
		xml_path,
		pdf_path,
		created_at,
		updated_at
	FROM sales_order_invoice
	WHERE sales_order_id = ?
	ORDER BY id DESC
	LIMIT 1
	`

	var inv model.SalesOrderInvoice

	var nfeKey sql.NullString
	var number sql.NullInt64
	var series sql.NullInt64
	var statusCode sql.NullString
	var statusReason sql.NullString
	var protocol sql.NullString
	var issuedAt sql.NullTime
	var authorizedAt sql.NullTime
	var cancelledAt sql.NullTime
	var xmlPath sql.NullString
	var pdfPath sql.NullString

	err := d.db.QueryRow(query, salesOrderID).Scan(
		&inv.ID,
		&inv.SalesOrderID,
		&nfeKey,
		&number,
		&series,
		&inv.Status,
		&statusCode,
		&statusReason,
		&protocol,
		&issuedAt,
		&authorizedAt,
		&cancelledAt,
		&xmlPath,
		&pdfPath,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if nfeKey.Valid {
		inv.NFeKey = &nfeKey.String
	}
	if number.Valid {
		inv.Number = &number.Int64
	}
	if series.Valid {
		v := int(series.Int64)
		inv.Series = &v
	}
	if statusCode.Valid {
		inv.StatusCode = &statusCode.String
	}
	if statusReason.Valid {
		inv.StatusReason = &statusReason.String
	}
	if protocol.Valid {
		inv.Protocol = &protocol.String
	}
	if issuedAt.Valid {
		inv.IssuedAt = &issuedAt.Time
	}
	if authorizedAt.Valid {
		inv.AuthorizedAt = &authorizedAt.Time
	}
	if cancelledAt.Valid {
		inv.CancelledAt = &cancelledAt.Time
	}
	if xmlPath.Valid {
		inv.XMLPath = &xmlPath.String
	}
	if pdfPath.Valid {
		inv.PDFPath = &pdfPath.String
	}

	return &inv, nil
}

func (d *salesOrderDAO) Update(s *model.SalesOrder) error {

	query := `
	UPDATE sales_order
	SET
		price = ?,
		discount = ?,
		status = ?,
		payment_method_id = ?,
		active = ?,

		customer_name = ?,
		customer_document = ?,

		marketplace_id = ?,
		external_order_id = ?,
		external_pack_id = ?,

		payment_status = ?,
		delivery_status = ?

	WHERE id = ?
	`

	_, err := d.db.Exec(
		query,

		s.Price.StringFixed(2),
		s.Discount.StringFixed(2),
		s.Status,
		s.PaymentMethodID,
		s.Active,

		s.CustomerName,
		s.CustomerDocument,

		s.MarketplaceID,
		s.ExternalOrderID,
		s.ExternalPackID,

		s.PaymentStatus,
		s.DeliveryStatus,

		s.ID,
	)

	return err
}

func (d *salesOrderDAO) Delete(id int64) error {
	_, err := d.db.Exec(`DELETE FROM sales_order WHERE id=?`, id)
	return err
}

func (d *salesOrderDAO) List() ([]model.SalesOrder, error) {

	query := `
	SELECT
		so.id,
		so.price,
		so.discount,
		so.status,
		so.payment_method_id,
		so.active,
		so.created_at,
		so.customer_name,
		so.customer_document,
		so.marketplace_id,
		so.external_order_id,
		so.external_pack_id,
		so.payment_status,
		so.delivery_status,
		so.updated_at,

		pm.id,
		pm.name,
		pm.active,

		m.id,
		m.name,
		m.description,
		m.commission_rate,
		m.active

	FROM sales_order so

	INNER JOIN payment_method pm
		ON pm.id = so.payment_method_id

	LEFT JOIN marketplace m
		ON m.id = so.marketplace_id

	ORDER BY so.created_at DESC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.SalesOrder

	for rows.Next() {

		var s model.SalesOrder

		var price, discount string

		var customerName sql.NullString
		var customerDocument sql.NullString

		var marketplaceID sql.NullInt64
		var externalOrderID sql.NullString
		var externalPackID sql.NullString

		var paymentStatus sql.NullString
		var deliveryStatus sql.NullString

		var updatedAt sql.NullTime

		var paymentMethod model.PaymentMethod

		var marketplace model.Marketplace
		var marketplaceIDJoined sql.NullInt64
		var marketplaceName sql.NullString
		var marketplaceDescription sql.NullString
		var marketplaceCommissionRate sql.NullString
		var marketplaceActive sql.NullBool

		err := rows.Scan(

			&s.ID,
			&price,
			&discount,
			&s.Status,
			&s.PaymentMethodID,
			&s.Active,
			&s.CreatedAt,

			&customerName,
			&customerDocument,

			&marketplaceID,
			&externalOrderID,
			&externalPackID,

			&paymentStatus,
			&deliveryStatus,

			&updatedAt,

			&paymentMethod.ID,
			&paymentMethod.Name,
			&paymentMethod.Active,

			&marketplaceIDJoined,
			&marketplaceName,
			&marketplaceDescription,
			&marketplaceCommissionRate,
			&marketplaceActive,
		)

		if err != nil {
			return nil, err
		}

		s.Price, _ = decimal.NewFromString(price)
		s.Discount, _ = decimal.NewFromString(discount)

		if customerName.Valid {
			s.CustomerName = &customerName.String
		}

		if customerDocument.Valid {
			s.CustomerDocument = &customerDocument.String
		}

		if marketplaceID.Valid {
			s.MarketplaceID = &marketplaceID.Int64
		}

		if externalOrderID.Valid {
			s.ExternalOrderID = &externalOrderID.String
		}

		if externalPackID.Valid {
			s.ExternalPackID = &externalPackID.String
		}

		if paymentStatus.Valid {
			s.PaymentStatus = &paymentStatus.String
		}

		if deliveryStatus.Valid {
			s.DeliveryStatus = &deliveryStatus.String
		}

		if updatedAt.Valid {
			s.UpdatedAt = &updatedAt.Time
		}

		s.PaymentMethod = &paymentMethod

		if marketplaceIDJoined.Valid {

			marketplace.ID = marketplaceIDJoined.Int64
			marketplace.Name = marketplaceName.String

			if marketplaceCommissionRate.Valid {
				marketplace.CommissionRate, _ =
					decimal.NewFromString(marketplaceCommissionRate.String)
			}

			s.Marketplace = &marketplace
		}

		invoice, err := d.GetInvoiceBySalesOrderID(s.ID)
		if err == nil {
			s.Invoice = invoice
		}

		list = append(list, s)
	}

	return list, nil
}

type SalesOrderItemDAO interface {
	Create(i *model.SalesOrderItem) error
	ListByOrder(orderID int64) ([]model.SalesOrderItem, error)
	DeleteByOrder(orderID int64) error
}

type salesOrderItemDAO struct {
	db *sql.DB
}

func NewSalesOrderItemDAO() SalesOrderItemDAO {
	return &salesOrderItemDAO{db: config.NewDatabaseConfig().DB}
}

func (d *salesOrderItemDAO) Create(i *model.SalesOrderItem) error {

	query := `
	INSERT INTO sales_order_item
	(
		sales_order_id,
		purchase_item_id,
		quantity,
		unit_price
	)
	VALUES (?, ?, ?, ?)
	`

	_, err := d.db.Exec(
		query,
		i.SalesOrderID,
		i.PurchaseItemID,
		i.Quantity,
		i.UnitPrice.StringFixed(2),
	)

	return err
}

func (d *salesOrderItemDAO) ListByOrder(orderID int64) ([]model.SalesOrderItem, error) {

	query := `
	SELECT
		soi.sales_order_id,
		soi.purchase_item_id,
		soi.quantity,
		soi.unit_price,

		pi.id,
		pi.purchase_id,
		pi.product_id,
		pi.quantity,
		pi.cost_price,
		pi.total,
		pi.cost_center_id

	FROM sales_order_item soi

	INNER JOIN purchase_item pi
		ON pi.id = soi.purchase_item_id

	WHERE soi.sales_order_id = ?
	`

	rows, err := d.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.SalesOrderItem

	for rows.Next() {

		var item model.SalesOrderItem

		var unitPrice string

		var purchaseItem model.PurchaseItem

		var purchaseQty string
		var costPrice string
		var purchaseTotal string

		var costCenterID sql.NullInt64

		err := rows.Scan(

			&item.SalesOrderID,
			&item.PurchaseItemID,
			&item.Quantity,
			&unitPrice,

			&purchaseItem.ID,
			&purchaseItem.PurchaseID,
			&purchaseItem.ProductID,
			&purchaseQty,
			&costPrice,
			&purchaseTotal,
			&costCenterID,
		)

		if err != nil {
			return nil, err
		}

		item.UnitPrice, _ =
			decimal.NewFromString(unitPrice)

		purchaseItem.Quantity, _ =
			decimal.NewFromString(purchaseQty)

		purchaseItem.CostPrice, _ =
			decimal.NewFromString(costPrice)

		purchaseItem.Total, _ =
			decimal.NewFromString(purchaseTotal)

		if costCenterID.Valid {
			purchaseItem.CostCenterID = &costCenterID.Int64
		}

		item.PurchaseItem = &purchaseItem

		list = append(list, item)
	}

	return list, nil
}

func (d *salesOrderItemDAO) DeleteByOrder(orderID int64) error {
	_, err := d.db.Exec(
		`DELETE FROM sales_order_item WHERE sales_order_id=?`,
		orderID,
	)
	return err
}
