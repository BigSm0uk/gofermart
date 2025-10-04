-- Удаление триггеров
DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Удаление функции
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаление индексов
DROP INDEX IF EXISTS idx_loyalty_operations_processed_at;
DROP INDEX IF EXISTS idx_loyalty_operations_order_number;
DROP INDEX IF EXISTS idx_loyalty_operations_user_id;
DROP INDEX IF EXISTS idx_orders_uploaded_at;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_number;
DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_users_login;

-- Удаление таблиц
DROP TABLE IF EXISTS loyalty_operations;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS users;
