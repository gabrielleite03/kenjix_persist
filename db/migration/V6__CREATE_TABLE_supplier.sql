CREATE TABLE supplier (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    country VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(100),
    active BOOLEAN DEFAULT true
);
