CREATE TABLE stock (
    product_id BIGINT NOT NULL,
    warehouse_id BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    active BOOLEAN DEFAULT true,

    PRIMARY KEY (product_id, warehouse_id),
    FOREIGN KEY (product_id) REFERENCES product(id),
    FOREIGN KEY (warehouse_id) REFERENCES warehouse(id)
);
