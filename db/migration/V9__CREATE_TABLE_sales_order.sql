CREATE TABLE payment_method (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    active BOOLEAN DEFAULT TRUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sales_order (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    price DECIMAL(12,2) NOT NULL,
    discount DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    status VARCHAR(30) NOT NULL,
    payment_method_id BIGINT NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_sales_order_payment_method
        FOREIGN KEY (payment_method_id)
        REFERENCES payment_method(id)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE sales_order
ADD COLUMN customer_name VARCHAR(150) NULL,
ADD COLUMN customer_document VARCHAR(20) NULL,

ADD COLUMN marketplace_id BIGINT NULL,
ADD COLUMN external_order_id VARCHAR(100) NULL,
ADD COLUMN external_pack_id VARCHAR(100) NULL,

ADD COLUMN payment_status VARCHAR(30) NULL,
ADD COLUMN delivery_status VARCHAR(30) NULL,

ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;

CREATE TABLE sales_order_invoice (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,

    sales_order_id BIGINT NOT NULL,

    nfe_key VARCHAR(44) NULL,
    number BIGINT NULL,
    series INT NULL,

    status VARCHAR(30) NOT NULL,
    status_code VARCHAR(10) NULL,
    status_reason VARCHAR(255) NULL,

    protocol VARCHAR(50) NULL,
    issued_at DATETIME NULL,
    authorized_at DATETIME NULL,
    cancelled_at DATETIME NULL,

    xml_path TEXT NULL,
    pdf_path TEXT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_so_invoice_order
        FOREIGN KEY (sales_order_id)
        REFERENCES sales_order(id),

    UNIQUE KEY uk_sales_order_invoice_key (nfe_key),
    INDEX idx_so_invoice_order (sales_order_id),
    INDEX idx_so_invoice_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
