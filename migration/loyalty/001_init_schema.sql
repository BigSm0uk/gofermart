-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE balances (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current NUMERIC(12,2) NOT NULL DEFAULT 0,
    withdrawn NUMERIC(12,2) NOT NULL DEFAULT 0
);

CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_number TEXT UNIQUE NOT NULL,
    status order_status NOT NULL DEFAULT 'NEW',
    accrual NUMERIC(12,2),
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE withdrawals (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_number TEXT NOT NULL,
    sum NUMERIC(12,2) NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS withdrawals;
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS order_status;
DROP TABLE IF EXISTS balances;
DROP TABLE IF EXISTS users;
