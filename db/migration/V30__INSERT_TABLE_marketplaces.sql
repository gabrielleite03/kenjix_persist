INSERT INTO marketplace (
    id, name, logo, comission_rate, integration_type, api_endpoint, active
) VALUES
(1, 'Mercado Livre', 
    'https://grandesnomesdapropaganda.com.br/wp-content/uploads/2021/04/Mercado-Livre.jpg',
    16, 'api', 'https://api.mercadolibre.com', 1
),
(2, 'Shopee',
    'https://logodownload.org/wp-content/uploads/2021/03/shopee-logo-0.png',
    18, 'api', 'https://partner.shopeemobile.com/api/v2', 1
),
(3, 'Amazon',
    'https://logodownload.org/wp-content/uploads/2014/04/amazon-logo-0.png',
    15, 'manual', NULL, 0