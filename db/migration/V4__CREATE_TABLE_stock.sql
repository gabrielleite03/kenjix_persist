CREATE TABLE stock (
    product_id BIGINT NOT NULL,
    warehouse_place_id BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT true,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (product_id, warehouse_place_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;