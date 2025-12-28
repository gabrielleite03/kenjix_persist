CREATE TABLE import_process (
    id BIGSERIAL PRIMARY KEY,
    purchase_order_id BIGINT UNIQUE NOT NULL,
    incoterm VARCHAR(10) NOT NULL,
    exchange_rate NUMERIC(10,4) NOT NULL,
    status VARCHAR(30) NOT NULL,
    arrival_date DATE,

    FOREIGN KEY (purchase_order_id) REFERENCES purchase_order(id)
);
