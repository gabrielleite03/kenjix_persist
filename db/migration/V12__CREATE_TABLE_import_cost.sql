CREATE TABLE import_cost (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    import_process_id BIGINT NOT NULL,
    type VARCHAR(30) NOT NULL, -- FREIGHT, TAX, INSURANCE...
    description VARCHAR(255),
    amount DECIMAL(14,2) NOT NULL,
    currency VARCHAR(10) NOT NULL,

    CONSTRAINT fk_import_cost_import_process
        FOREIGN KEY (import_process_id)
        REFERENCES import_process(id)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;