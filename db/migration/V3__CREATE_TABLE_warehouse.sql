CREATE TABLE warehouse (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    active BOOLEAN DEFAULT true
);
