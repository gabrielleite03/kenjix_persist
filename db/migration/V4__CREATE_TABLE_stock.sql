CREATE TABLE stock (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,

    product_id BIGINT NOT NULL,
    warehouse_place_id BIGINT NOT NULL,
    purchase_item_id BIGINT NOT NULL,

    quantity INT NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT TRUE,

    updated_at TIMESTAMP 
        DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uq_stock_compound (
        product_id,
        warehouse_place_id,
        purchase_item_id
    ),

    INDEX idx_stock_product (product_id),
    INDEX idx_stock_place (warehouse_place_id),
    INDEX idx_stock_purchase_item (purchase_item_id)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE stock 
ADD CONSTRAINT fk_stock_product
FOREIGN KEY (product_id) REFERENCES product(id);

ALTER TABLE stock 
ADD CONSTRAINT fk_stock_place
FOREIGN KEY (warehouse_place_id) REFERENCES warehouse_place(id);

ALTER TABLE stock 
ADD CONSTRAINT fk_stock_purchase_item
FOREIGN KEY (purchase_item_id) REFERENCES purchase_item(id);