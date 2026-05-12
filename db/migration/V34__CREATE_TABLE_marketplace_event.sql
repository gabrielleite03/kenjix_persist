CREATE TABLE marketplace_event (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    source VARCHAR(50) NOT NULL,
    topic VARCHAR(80) NOT NULL,
    resource VARCHAR(255) NOT NULL,
    user_id BIGINT NULL,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    error_message TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at DATETIME NULL,

    UNIQUE KEY uk_marketplace_event (source, topic, resource)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;