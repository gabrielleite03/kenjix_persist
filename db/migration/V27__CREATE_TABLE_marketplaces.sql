CREATE TABLE marketplace (
    id BIGINT NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    logo VARCHAR(500),
    comission_rate DECIMAL(10,2) NOT NULL DEFAULT 0.0,
    integration_type VARCHAR(100),
    api_key VARCHAR(255),
    api_secret VARCHAR(255),
    api_endpoint VARCHAR(500),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (id)
);