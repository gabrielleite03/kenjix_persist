CREATE INDEX idx_product_sku ON product(sku);
CREATE INDEX idx_stock_product ON stock(product_id);
CREATE INDEX idx_stock_movement_created ON stock_movement(created_at);
CREATE INDEX idx_purchase_order_status ON purchase_order(status);
CREATE INDEX idx_import_process_status ON import_process(status);
