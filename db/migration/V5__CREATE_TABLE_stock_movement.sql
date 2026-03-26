CREATE TABLE stock_movement (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    product_id BIGINT NOT NULL,
    warehouse_id BIGINT NOT NULL,
    type ENUM('IN','OUT','ADJUSTMENT') NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reason VARCHAR(20) NOT NULL,

    CONSTRAINT fk_stock_movement_product
        FOREIGN KEY (product_id) REFERENCES product(id),

    CONSTRAINT fk_stock_movement_warehouse
        FOREIGN KEY (warehouse_id) REFERENCES warehouse(id)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;