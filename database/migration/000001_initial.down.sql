-- =========================================================
-- Rollback Tailor Order Management Schema (Initial)
-- PostgreSQL
-- =========================================================

-- Drop update triggers first.
DROP TRIGGER IF EXISTS trg_notification_logs_updated_at ON notification_logs;
DROP TRIGGER IF EXISTS trg_order_items_updated_at ON order_items;
DROP TRIGGER IF EXISTS trg_orders_updated_at ON orders;
DROP TRIGGER IF EXISTS trg_customers_updated_at ON customers;
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;

-- Drop tables in reverse dependency order.
DROP TABLE IF EXISTS notification_logs;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS service_types;
DROP TABLE IF EXISTS clothes_categories;
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS users;

-- Drop helper function after all triggers/tables are removed.
DROP FUNCTION IF EXISTS set_updated_at();

-- Drop extension created by initial migration.
DROP EXTENSION IF EXISTS citext;
