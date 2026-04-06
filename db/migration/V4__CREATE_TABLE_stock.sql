CREATE TABLE stock (
    product_id BIGINT NOT NULL,
    warehouse_id BIGINT NOT NULL,
    warehouse_place_id BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    active BOOLEAN DEFAULT TRUE,

    PRIMARY KEY (product_id, warehouse_id),

    CONSTRAINT fk_stock_product
        FOREIGN KEY (product_id) REFERENCES product(id),

    CONSTRAINT fk_stock_warehouse
        FOREIGN KEY (warehouse_id) REFERENCES warehouse(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;