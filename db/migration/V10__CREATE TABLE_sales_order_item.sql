CREATE TABLE sales_order_item (
    sales_order_id BIGINT NOT NULL,
    purchase_item_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(12,2) NOT NULL,

    PRIMARY KEY (sales_order_id, purchase_item_id),

    CONSTRAINT fk_soi_sales_order
        FOREIGN KEY (sales_order_id)
        REFERENCES sales_order(id),

    CONSTRAINT fk_soi_product
        FOREIGN KEY (purchase_item_id)
        REFERENCES purchase_item(id)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE sales_order_item
DROP PRIMARY KEY;

ALTER TABLE sales_order_item
ADD COLUMN id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY FIRST,
ADD COLUMN external_item_id VARCHAR(80) NULL,
ADD COLUMN external_sku VARCHAR(100) NULL;

CREATE INDEX idx_soi_order ON sales_order_item (sales_order_id);
CREATE INDEX idx_soi_product ON sales_order_item (purchase_item_id);