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
