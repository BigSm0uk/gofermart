-- +goose Up
-- Включение расширения для UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создание таблицы пользователей
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    login VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Создание таблицы заказов
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    number VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'NEW',
    accrual DECIMAL(10,2) DEFAULT 0,
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Создание таблицы операций лояльности
CREATE TABLE loyalty_operations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_number VARCHAR(255),
    operation_type VARCHAR(20) NOT NULL, -- 'CREDIT' или 'DEBIT'
    amount DECIMAL(10,2) NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Создание индексов для оптимизации запросов
CREATE INDEX idx_users_login ON users(login);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_number ON orders(number);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_uploaded_at ON orders(uploaded_at);
CREATE INDEX idx_loyalty_operations_user_id ON loyalty_operations(user_id);
CREATE INDEX idx_loyalty_operations_order_number ON loyalty_operations(order_number);
CREATE INDEX idx_loyalty_operations_processed_at ON loyalty_operations(processed_at);

-- Триггеры для автоматического обновления updated_at будут добавлены позже

-- +goose Down

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