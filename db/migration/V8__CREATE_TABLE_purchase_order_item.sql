CREATE TABLE purchase_item (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    purchase_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity DECIMAL(15,4) NOT NULL,
    cost_price DECIMAL(15,4) NOT NULL,
    total DECIMAL(15,2) NOT NULL,
    cost_center_id BIGINT NULL,

    CONSTRAINT fk_purchase_item_purchase
        FOREIGN KEY (purchase_id) REFERENCES purchase(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_purchase_item_product
        FOREIGN KEY (product_id) REFERENCES product(id),

    CONSTRAINT fk_purchase_item_cost_center
        FOREIGN KEY (cost_center_id) REFERENCES cost_center(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;