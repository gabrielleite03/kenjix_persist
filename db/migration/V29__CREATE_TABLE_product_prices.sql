CREATE TABLE `product_prices` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `product_id` BIGINT NOT NULL,
    `marketplace_id` BIGINT NOT NULL,
    `price` DECIMAL(14,2) NOT NULL,
    `active` TINYINT(1) NOT NULL DEFAULT 1,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_product_marketplace` (`product_id`, `marketplace_id`),
    CONSTRAINT `fk_product_prices_product` FOREIGN KEY (`product_id`)
        REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_product_prices_marketplace` FOREIGN KEY (`marketplace_id`)
        REFERENCES `marketplace` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;