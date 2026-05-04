BEGIN;

-- Migrate order_items table
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS name_th text;
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS name_en text;

UPDATE order_items SET name_th = name WHERE name_th IS NULL OR name_th = '';
ALTER TABLE order_items ALTER COLUMN name_th SET NOT NULL;

COMMIT;
