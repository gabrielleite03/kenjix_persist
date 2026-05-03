CREATE TABLE product_marketplace (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,

    product_id BIGINT NOT NULL,
    marketplace_id BIGINT NOT NULL,

    external_id VARCHAR(100) NULL, -- ID do produto no marketplace (ex: MLB123456789)
    product_url TEXT NOT NULL,     -- link do produto no marketplace

    price DECIMAL(12,2) NULL,      -- preço naquele marketplace
    listing_type VARCHAR(50) NULL, -- gold_pro, gold_special etc
    status VARCHAR(20) NULL,       -- active, paused, closed

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    active TINYINT(1) NOT NULL DEFAULT 1,

    CONSTRAINT fk_pm_product
        FOREIGN KEY (product_id) REFERENCES product(id),

    CONSTRAINT fk_pm_marketplace
        FOREIGN KEY (marketplace_id) REFERENCES marketplace(id),

    UNIQUE KEY uniq_product_marketplace (product_id, marketplace_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;