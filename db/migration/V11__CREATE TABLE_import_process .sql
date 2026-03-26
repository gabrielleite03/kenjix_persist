CREATE TABLE import_process (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    purchase_order_id BIGINT NOT NULL,
    incoterm VARCHAR(10) NOT NULL,
    exchange_rate DECIMAL(10,4) NOT NULL,
    status VARCHAR(30) NOT NULL,
    arrival_date DATE,

    CONSTRAINT uk_import_process_po UNIQUE (purchase_order_id),

    CONSTRAINT fk_import_process_purchase_order
        FOREIGN KEY (purchase_order_id)
        REFERENCES purchase_order(id)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
