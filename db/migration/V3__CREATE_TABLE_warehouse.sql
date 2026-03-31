CREATE TABLE warehouse (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(300) NOT NULL,
    capacity BIGINT,
    active BOOLEAN DEFAULT TRUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE warehouse_place_type (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    active BOOLEAN DEFAULT TRUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE warehouse_place (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    capacity BIGINT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    warehouse_place_type_id BIGINT,
    warehouse_id BIGINT,
    CONSTRAINT fk_warehouse_place_type
        FOREIGN KEY (warehouse_place_type_id)
        REFERENCES warehouse_place_type(id),
    CONSTRAINT fk_warehouse_place_warehouse
        FOREIGN KEY (warehouse_id)
        REFERENCES warehouse(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
