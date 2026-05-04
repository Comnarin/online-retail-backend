-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Tenants
CREATE TABLE IF NOT EXISTS tenants (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug        VARCHAR(100) UNIQUE NOT NULL,
    name        VARCHAR(255) NOT NULL,
    liff_id     VARCHAR(255) DEFAULT 'PLACEHOLDER_LIFF_ID',
    features    JSONB NOT NULL DEFAULT '{"enable_membership":false,"enable_coupon":false,"enable_points":false,"enable_reviews":false,"enable_delivery":false}',
    appearance  JSONB NOT NULL DEFAULT '{"primary_color":"#06C755","secondary_color":"#004c1b","logo_url":"","banner_url":"","font_family":"Plus Jakarta Sans","store_name":"My Store","store_tagline":""}',
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Users
CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID REFERENCES tenants(id) ON DELETE SET NULL,
    line_user_id  VARCHAR(255) UNIQUE,
    email         VARCHAR(255) UNIQUE,
    password      VARCHAR(255),
    name          VARCHAR(255) NOT NULL,
    phone         VARCHAR(30),
    avatar        TEXT,
    role          VARCHAR(30) NOT NULL DEFAULT 'customer' CHECK (role IN ('superadmin','tenant_admin','customer')),
    is_active     BOOLEAN NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Product Categories
CREATE TABLE IF NOT EXISTS product_categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    parent_id   UUID REFERENCES product_categories(id) ON DELETE SET NULL,
    name        VARCHAR(255) NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_categories_tenant ON product_categories(tenant_id);

-- Products
CREATE TABLE IF NOT EXISTS products (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    category_id UUID REFERENCES product_categories(id) ON DELETE SET NULL,
    sku         VARCHAR(100),
    name        VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    price       NUMERIC(12,2) NOT NULL,
    stock       INT NOT NULL DEFAULT 0,
    images      JSONB NOT NULL DEFAULT '[]',
    status      VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('active','inactive','draft')),
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_tenant ON products(tenant_id);
CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);

-- Orders
CREATE TABLE IF NOT EXISTS orders (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id      UUID NOT NULL REFERENCES users(id),
    order_number     VARCHAR(50) UNIQUE NOT NULL,
    status           VARCHAR(30) NOT NULL DEFAULT 'pending'
                         CHECK (status IN ('pending','confirmed','processing','shipping','delivered','cancelled')),
    subtotal         NUMERIC(12,2) NOT NULL,
    discount_amount  NUMERIC(12,2) NOT NULL DEFAULT 0,
    total            NUMERIC(12,2) NOT NULL,
    coupon_id        UUID,
    points_used      INT NOT NULL DEFAULT 0,
    shipping_address JSONB NOT NULL DEFAULT '{}',
    note             TEXT DEFAULT '',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_tenant ON orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_orders_customer ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);

-- Order Items
CREATE TABLE IF NOT EXISTS order_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id    UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id  UUID NOT NULL REFERENCES products(id),
    name        VARCHAR(255) NOT NULL,
    price       NUMERIC(12,2) NOT NULL,
    quantity    INT NOT NULL,
    subtotal    NUMERIC(12,2) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);

-- Membership Tiers
CREATE TABLE IF NOT EXISTS membership_tiers (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name          VARCHAR(100) NOT NULL,
    min_points    INT NOT NULL DEFAULT 0,
    discount_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    color         VARCHAR(20) DEFAULT '#06C755',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Customer Memberships
CREATE TABLE IF NOT EXISTS customer_memberships (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tier_id    UUID NOT NULL REFERENCES membership_tiers(id),
    joined_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    UNIQUE(tenant_id, user_id)
);

-- Coupons
CREATE TABLE IF NOT EXISTS coupons (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    code            VARCHAR(50) NOT NULL,
    description     TEXT DEFAULT '',
    discount_type   VARCHAR(20) NOT NULL DEFAULT 'percent' CHECK (discount_type IN ('percent','fixed')),
    discount_value  NUMERIC(12,2) NOT NULL,
    min_order_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    max_uses        INT,
    used_count      INT NOT NULL DEFAULT 0,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    starts_at       TIMESTAMPTZ,
    expires_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

-- Point Transactions
CREATE TABLE IF NOT EXISTS point_transactions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_id    UUID REFERENCES orders(id) ON DELETE SET NULL,
    points      INT NOT NULL, -- positive=earn, negative=redeem
    description VARCHAR(255) DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_points_tenant_user ON point_transactions(tenant_id, user_id);

-- Seed: default superadmin
INSERT INTO users (id, name, email, password, role)
VALUES (
    gen_random_uuid(),
    'Super Admin',
    'admin@retail.com',
    -- Default password: 'Admin@1234' (bcrypt) -- CHANGE IN PRODUCTION
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'superadmin'
) ON CONFLICT (email) DO NOTHING;
