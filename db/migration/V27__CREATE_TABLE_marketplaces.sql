CREATE TABLE marketplace (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    logo TEXT NULL,

    status VARCHAR(20) NOT NULL,
    commission_rate DECIMAL(10,4) NOT NULL,
    integration_type VARCHAR(20) NOT NULL,

    api_url TEXT NULL,
    api_key TEXT NULL,
    api_secret TEXT NULL,
    api_endpoint TEXT NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,

    INDEX idx_marketplace_status (status),
    INDEX idx_marketplace_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;