CREATE TABLE purchase_order (
    id BIGSERIAL PRIMARY KEY,
    supplier_id BIGINT NOT NULL,
    status VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (supplier_id) REFERENCES supplier(id)
);
