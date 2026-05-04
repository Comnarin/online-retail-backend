BEGIN;

-- Migrate product_categories table
ALTER TABLE product_categories ADD COLUMN IF NOT EXISTS name_th text;
ALTER TABLE product_categories ADD COLUMN IF NOT EXISTS name_en text;

UPDATE product_categories SET name_th = name WHERE name_th IS NULL OR name_th = '';
ALTER TABLE product_categories ALTER COLUMN name_th SET NOT NULL;

-- Migrate products table
ALTER TABLE products ADD COLUMN IF NOT EXISTS name_th text;
ALTER TABLE products ADD COLUMN IF NOT EXISTS name_en text;
ALTER TABLE products ADD COLUMN IF NOT EXISTS description_th text;
ALTER TABLE products ADD COLUMN IF NOT EXISTS description_en text;

UPDATE products SET name_th = name WHERE name_th IS NULL OR name_th = '';
ALTER TABLE products ALTER COLUMN name_th SET NOT NULL;

UPDATE products SET description_th = description WHERE description_th IS NULL OR description_th = '';

COMMIT;
