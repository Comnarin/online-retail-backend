-- Add QR code URL to payment methods (admin uploads QR code for each method)
ALTER TABLE payment_methods ADD COLUMN IF NOT EXISTS qr_code_url TEXT DEFAULT '';

-- Add slip image URL to orders (customer uploads payment slip)
ALTER TABLE orders ADD COLUMN IF NOT EXISTS slip_image_url TEXT DEFAULT '';

-- Update order status check constraint to include new statuses
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_status_check;
ALTER TABLE orders ADD CONSTRAINT orders_status_check 
    CHECK (status IN ('pending','pending_verification','confirmed','processing','shipping','delivered','cancelled','expired'));

-- Increase default expiry for payment methods (manual QR flow needs more time)
ALTER TABLE payment_methods ALTER COLUMN expiry_minutes SET DEFAULT 30;
