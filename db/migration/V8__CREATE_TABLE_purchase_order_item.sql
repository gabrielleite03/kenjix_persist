CREATE TABLE purchase_order_item (
    purchase_order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(12,2) NOT NULL,
    active BOOLEAN DEFAULT TRUE,

    PRIMARY KEY (purchase_order_id, product_id),

    CONSTRAINT fk_poi_purchase_order
        FOREIGN KEY (purchase_order_id)
        REFERENCES purchase_order(id),

    CONSTRAINT fk_poi_product
        FOREIGN KEY (product_id)
        REFERENCES product(id)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;