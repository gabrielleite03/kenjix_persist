CREATE TABLE product_property (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    product_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    CONSTRAINT fk_product_property_product
        FOREIGN KEY (product_id) REFERENCES product(id)
);

CREATE TABLE product_image (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    product_id BIGINT NOT NULL,
    url TEXT NOT NULL,
    position INT DEFAULT 0,
    is_primary BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_product_image_product
        FOREIGN KEY (product_id) REFERENCES product(id)
);

CREATE TABLE product_video (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    product_id BIGINT NOT NULL,
    url TEXT NOT NULL,
    provider TEXT,
    CONSTRAINT fk_product_video_product
        FOREIGN KEY (product_id) REFERENCES product(id)
);