CREATE TABLE supplier (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    
    razao_social VARCHAR(255) NOT NULL,
    nome_fantasia VARCHAR(255) NOT NULL,
    cnpj VARCHAR(20) NOT NULL,
    
    ie VARCHAR(50) NULL,
    address VARCHAR(255) NULL,
    sales_person VARCHAR(255) NULL,
    email VARCHAR(255) NULL,
    phone VARCHAR(50) NULL,
    
    active BOOLEAN NOT NULL DEFAULT TRUE,
    
    category_id BIGINT NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_supplier_category
        FOREIGN KEY (category_id)
        REFERENCES category(id)
);
