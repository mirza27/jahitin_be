-- =========================================================
-- Tailor Order Management Schema
-- PostgreSQL
-- =========================================================

-- Optional: for case-insensitive unique email/username
CREATE EXTENSION IF NOT EXISTS citext;

-- =========================================================
-- USERS
-- =========================================================
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    username CITEXT UNIQUE,
    email CITEXT UNIQUE,
    password_hash TEXT,
    device_id VARCHAR(255) UNIQUE,
    user_type VARCHAR(30) NOT NULL DEFAULT 'tailor',
    phone VARCHAR(30),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =========================================================
-- CUSTOMERS
-- Satu user/penjahit punya banyak customer
-- =========================================================
CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(30),
    notes TEXT,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_customers_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_customers_user_id ON customers(user_id);
CREATE INDEX idx_customers_name ON customers(name);

-- =========================================================
-- CLOTHES CATEGORIES
-- Contoh: Gamis, Kebaya, Rok, Celana
-- =========================================================
CREATE TABLE clothes_categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =========================================================
-- SERVICE TYPES
-- Contoh: Jahit baru, Kecilkan, Besarkan, Ganti resleting
-- =========================================================
CREATE TABLE service_types (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =========================================================
-- ORDERS
-- Satu order milik satu user dan satu customer
-- =========================================================
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    user_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,

    deadline TIMESTAMPTZ,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    sharing_code VARCHAR(50) UNIQUE,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ,

    CONSTRAINT fk_orders_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_orders_customer
        FOREIGN KEY (customer_id)
        REFERENCES customers(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_deadline ON orders(deadline);
CREATE INDEX idx_orders_status ON orders(status);

-- =========================================================
-- ORDER ITEMS
-- Unit terkecil: pekerjaan jahit
-- Bisa jahit baru, permak, kecilkan, dll
-- =========================================================
CREATE TABLE order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    category_id BIGINT,
    service_type_id BIGINT,
    clothes_for VARCHAR(100) NOT NULL,
    custom_service_name VARCHAR(100),
    notes TEXT,
    price BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    finished_at TIMESTAMPTZ,

    CONSTRAINT fk_order_items_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_order_items_category
        FOREIGN KEY (category_id)
        REFERENCES clothes_categories(id)
        ON DELETE SET NULL,

    CONSTRAINT fk_order_items_service_type
        FOREIGN KEY (service_type_id)
        REFERENCES service_types(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_category_id ON order_items(category_id);
CREATE INDEX idx_order_items_service_type_id ON order_items(service_type_id);
CREATE INDEX idx_order_items_status ON order_items(status);

-- =========================================================
-- NOTIFICATION LOGS
-- Log pengingat deadline / status
-- =========================================================
CREATE TABLE notification_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    order_id BIGINT,
    order_item_id BIGINT,
    device_id VARCHAR(255),
    notification_type VARCHAR(50) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'sent',
    message TEXT,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_notification_logs_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_notification_logs_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE SET NULL,

    CONSTRAINT fk_notification_logs_order_item
        FOREIGN KEY (order_item_id)
        REFERENCES order_items(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_notification_logs_user_id ON notification_logs(user_id);
CREATE INDEX idx_notification_logs_order_id ON notification_logs(order_id);
CREATE INDEX idx_notification_logs_order_item_id ON notification_logs(order_item_id);

-- =========================================================
-- UPDATED_AT trigger helper
-- =========================================================
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_customers_updated_at
BEFORE UPDATE ON customers
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_order_items_updated_at
BEFORE UPDATE ON order_items
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_notification_logs_updated_at
BEFORE UPDATE ON notification_logs
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();