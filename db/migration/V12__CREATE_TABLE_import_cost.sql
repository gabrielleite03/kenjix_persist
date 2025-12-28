CREATE TABLE import_cost (
    id BIGSERIAL PRIMARY KEY,
    import_process_id BIGINT NOT NULL,
    type VARCHAR(30) NOT NULL, -- FREIGHT, TAX, INSURANCE...
    description VARCHAR(255),
    amount NUMERIC(14,2) NOT NULL,
    currency VARCHAR(10) NOT NULL,

    FOREIGN KEY (import_process_id) REFERENCES import_process(id)
);
