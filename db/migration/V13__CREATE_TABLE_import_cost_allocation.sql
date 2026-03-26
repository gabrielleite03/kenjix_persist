CREATE TABLE import_cost_allocation (
    import_cost_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    allocated_amount DECIMAL(14,2) NOT NULL,

    PRIMARY KEY (import_cost_id, product_id),

    CONSTRAINT fk_import_cost_allocation_cost
        FOREIGN KEY (import_cost_id)
        REFERENCES import_cost(id),

    CONSTRAINT fk_import_cost_allocation_product
        FOREIGN KEY (product_id)
        REFERENCES product(id)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
