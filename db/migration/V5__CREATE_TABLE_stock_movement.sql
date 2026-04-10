CREATE TABLE stock_movement (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,

    product_id BIGINT NOT NULL,
    warehouse_place_id BIGINT NOT NULL,
    purchase_item_id BIGINT NOT NULL,

    type ENUM('IN', 'OUT', 'ADJUSTMENT') NOT NULL,
    quantity INT NOT NULL,

    reference_id BIGINT NULL,
    reference_type VARCHAR(50) NULL,

    reason VARCHAR(255) NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_stock_movement_product (product_id),
    INDEX idx_stock_movement_place (warehouse_place_id),
    INDEX idx_stock_movement_purchase_item (purchase_item_id),
    INDEX idx_stock_movement_created (created_at),
    INDEX idx_stock_movement_reference (reference_id)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;