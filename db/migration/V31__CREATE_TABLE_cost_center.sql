CREATE TABLE cost_center (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) NOT NULL,
    description TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE cost_center_property (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    cost_center_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    value DECIMAL(14,2) NOT NULL,
    type ENUM('index', 'value') NOT NULL,

    CONSTRAINT fk_cost_center_property
        FOREIGN KEY (cost_center_id)
        REFERENCES cost_center(id)
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;