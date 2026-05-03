CREATE TABLE product (
    id BIGINT NOT NULL AUTO_INCREMENT,
    name VARCHAR(150) NOT NULL,
    sku VARCHAR(50) NOT NULL UNIQUE,
    price DECIMAL(12,2) NOT NULL,
    marca VARCHAR(150) NOT NULL,
    description VARCHAR(500) NOT NULL,
    active TINYINT(1) NOT NULL DEFAULT 1,
    volume DECIMAL(12,2) NOT NULL,
    category_id BIGINT,
    PRIMARY KEY (id),
    CONSTRAINT fk_product_category
        FOREIGN KEY (category_id)
        REFERENCES category(id)
) ENGINE=InnoDB;


ALTER TABLE product
ADD COLUMN ncm VARCHAR(8) NULL,
ADD COLUMN ean VARCHAR(14) NULL;

ALTER TABLE product
ADD COLUMN weight DECIMAL(12,2) NULL;

ALTER TABLE product
ADD COLUMN is_kit TINYINT(1) DEFAULT 0;


CREATE TABLE product_kit (
    id BIGINT NOT NULL AUTO_INCREMENT,
    product_id BIGINT NOT NULL,         -- o KIT (produto pai)
    component_product_id BIGINT NOT NULL, -- produto que compõe o kit
    quantity INT NOT NULL DEFAULT 1,

    PRIMARY KEY (id),

    CONSTRAINT fk_product_kit_parent
        FOREIGN KEY (product_id)
        REFERENCES product(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_product_kit_component
        FOREIGN KEY (component_product_id)
        REFERENCES product(id)
);