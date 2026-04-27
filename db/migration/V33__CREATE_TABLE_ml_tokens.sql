CREATE TABLE ml_tokens (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,

    user_id BIGINT NOT NULL, -- ID do usuário no Mercado Livre
    nickname VARCHAR(100),   -- opcional (ajuda debug)

    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,

    expires_at DATETIME NOT NULL,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_ml_user_id (user_id)
);