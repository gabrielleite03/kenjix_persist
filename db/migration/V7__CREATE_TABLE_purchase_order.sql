CREATE TABLE purchase_status (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    description VARCHAR(255),
    active BOOLEAN DEFAULT TRUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE fiscal_number_type (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    description VARCHAR(255),
    active BOOLEAN DEFAULT TRUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE purchase_order (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    purchase_status_id BIGINT NOT NULL,
    fiscal_number_type_id BIGINT NOT NULL,
    supplier_id BIGINT NOT NULL,
    fiscal_number VARCHAR(30),
    status VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active BOOLEAN DEFAULT TRUE,

    CONSTRAINT fk_purchase_order_supplier
        FOREIGN KEY (supplier_id) REFERENCES supplier(id),

    CONSTRAINT fk_purchase_order_purchase_status
        FOREIGN KEY (purchase_status_id) REFERENCES purchase_status(id),

    CONSTRAINT fk_purchase_order_fiscal_number_type
        FOREIGN KEY (fiscal_number_type_id) REFERENCES fiscal_number_type(id)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;