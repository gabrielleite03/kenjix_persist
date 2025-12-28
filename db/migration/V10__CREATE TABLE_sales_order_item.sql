CREATE TABLE sales_order_item (
    sales_order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    unit_price NUMERIC(12,2) NOT NULL,

    PRIMARY KEY (sales_order_id, product_id),
    FOREIGN KEY (sales_order_id) REFERENCES sales_order(id),
    FOREIGN KEY (product_id) REFERENCES product(id)
);
