CREATE TABLE purchase (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    invoice_number VARCHAR(50) NULL,
    invoice_type VARCHAR(20) NOT NULL,
    supplier_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL,
    total DECIMAL(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_purchase_supplier
        FOREIGN KEY (supplier_id) REFERENCES supplier(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;