CREATE TABLE supplier (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    cnpj VARCHAR(20) NOT NULL,
    country VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(100),
    seller VARCHAR(100),
    seller_fone VARCHAR(100),
    active BOOLEAN DEFAULT TRUE,
    category_id BIGINT,

    CONSTRAINT uk_supplier_cnpj UNIQUE (cnpj),

    CONSTRAINT fk_supplier_category
        FOREIGN KEY (category_id)
        REFERENCES category(id)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
