CREATE TABLE sales_order_item (
    sales_order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(12,2) NOT NULL,

    PRIMARY KEY (sales_order_id, product_id),

    CONSTRAINT fk_soi_sales_order
        FOREIGN KEY (sales_order_id)
        REFERENCES sales_order(id),

    CONSTRAINT fk_soi_product
        FOREIGN KEY (product_id)
        REFERENCES product(id)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;