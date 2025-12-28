CREATE TABLE import_cost_allocation (
    import_cost_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    allocated_amount NUMERIC(14,2) NOT NULL,

    PRIMARY KEY (import_cost_id, product_id),
    FOREIGN KEY (import_cost_id) REFERENCES import_cost(id),
    FOREIGN KEY (product_id) REFERENCES product(id)
);
