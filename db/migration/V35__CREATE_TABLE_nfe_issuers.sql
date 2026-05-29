CREATE TABLE nfe_issuers (
    id BIGINT NOT NULL AUTO_INCREMENT,

    cnpj VARCHAR(14) NULL,
    cpf VARCHAR(11) NULL,
    razao_social VARCHAR(255) NOT NULL,
    nome_fantasia VARCHAR(255) NULL,

    inscricao_estadual VARCHAR(20) NULL,
    crt VARCHAR(1) NOT NULL,

    logradouro VARCHAR(255) NOT NULL,
    numero VARCHAR(20) NOT NULL,
    complemento VARCHAR(255) NULL,
    bairro VARCHAR(120) NOT NULL,
    codigo_mun VARCHAR(7) NOT NULL,
    municipio VARCHAR(120) NOT NULL,
    uf CHAR(2) NOT NULL,
    cep VARCHAR(8) NOT NULL,
    codigo_pais VARCHAR(4) NOT NULL DEFAULT '1058',
    pais VARCHAR(60) NOT NULL DEFAULT 'Brasil',
    telefone VARCHAR(20) NULL,

    active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    UNIQUE KEY uk_nfe_issuers_cnpj (cnpj),
    UNIQUE KEY uk_nfe_issuers_cpf (cpf),

    CONSTRAINT chk_nfe_issuers_document
        CHECK (
            (cnpj IS NOT NULL AND cpf IS NULL)
            OR
            (cnpj IS NULL AND cpf IS NOT NULL)
        )
);



INSERT INTO nfe_issuers (
    cnpj,
    cpf,
    razao_social,
    nome_fantasia,
    inscricao_estadual,
    crt,
    logradouro,
    numero,
    complemento,
    bairro,
    codigo_mun,
    municipio,
    uf,
    cep,
    codigo_pais,
    pais,
    telefone,
    active
) VALUES (
    '65468523000102',
    NULL,
    'KENJI IMPORTACAO E COMERCIO LTDA',
    NULL,
    '158447676112',
    '1',
    'Rua A',
    '100',
    NULL,
    'Centro',
    '3550308',
    'Sao Paulo',
    'SP',
    '01001000',
    '1058',
    'Brasil',
    NULL,
    TRUE
);