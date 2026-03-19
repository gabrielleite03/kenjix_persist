@startuml
' =========================
' CATÁLOGO
' =========================
entity Category {
  +id : Long
  --
  name : String
  active: Boolean
}

entity Product {
  +id : Long
  --
  name : String
  sku : String
  price : Decimal
  active : Boolean
}

entity product_property {
  +id : BIGINT
  --
  product_id : BIGINT
  name : TEXT
  value : TEXT
}

entity product_image {
  +id : BIGINT
  --
  product_id : BIGINT
  url : TEXT
  position : INT
  is_primary : BOOLEAN
}

entity product_video {
  +id : BIGINT
  --
  product_id : BIGINT
  url : TEXT
  provider : TEXT
}

Category ||--o{ Product : "1:N"
Product ||--o{ product_property : "1:N"
Product ||--o{ product_image : "1:N"
Product ||--o{ product_video : "1:N"

' =========================
' LOCALIZAÇÃO / ESTOQUE
' =========================
entity Warehouse {
  +id : Long
  --
  name : String
  active: Boolean
}

entity Stock {
  +product_id : Long
  +warehouse_id : Long
  --
  quantity : Int
  active: Boolean
}

entity StockMovement {
  +id : Long
  --
  type : String
  quantity : Int
  created_at : DateTime
}

Product ||--o{ Stock
Warehouse ||--o{ Stock

Product ||--o{ StockMovement
Warehouse ||--o{ StockMovement

' =========================
' FORNECEDOR E COMPRAS
' =========================
entity Supplier {
  +id : Long
  --
  name : String
  country : String
  email : String
  phone : String
  active: Boolean
}

entity PurchaseOrder {
  +id : Long
  --
  status : String
  created_at : DateTime
  active: Boolean
}

entity PurchaseOrderItem {
  +purchase_order_id : Long
  +product_id : Long
  --
  quantity : Int
  unit_price : Decimal
}

Supplier ||--o{ PurchaseOrder
PurchaseOrder ||--o{ PurchaseOrderItem
Product ||--o{ PurchaseOrderItem

' =========================
' VENDAS
' =========================
entity SalesOrder {
  +id : Long
  --
  status : String
  created_at : DateTime
  active: Boolean
}

entity SalesOrderItem {
  +sales_order_id : Long
  +product_id : Long
  --
  quantity : Int
  unit_price : Decimal
}

SalesOrder ||--o{ SalesOrderItem
Product ||--o{ SalesOrderItem

' =========================
' IMPORTAÇÃO
' =========================
entity ImportProcess {
  +id : Long
  --
  incoterm : String
  exchange_rate : Decimal
  status : String
  arrival_date : Date
  active: Boolean
}

entity ImportCost {
  +id : Long
  --
  type : String
  description : String
  amount : Decimal
  currency : String
}

entity ImportCostAllocation {
  +import_cost_id : Long
  +product_id : Long
  --
  allocated_amount : Decimal
}

PurchaseOrder ||--|| ImportProcess
ImportProcess ||--o{ ImportCost
ImportCost ||--o{ ImportCostAllocation
Product ||--o{ ImportCostAllocation

@enduml

