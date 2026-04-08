INSERT INTO cost_center (name, code, description)
VALUES ('PDV', 'PDV', 'Ponto de venda');

INSERT INTO cost_center_property (cost_center_id, name, value, type)
VALUES 
('1', 'Margem', 150, 'index'),
('1', 'Frete', 1, 'value');