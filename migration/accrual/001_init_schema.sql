-- +goose Up
CREATE TYPE accrual_order_status AS ENUM ('REGISTERED', 'INVALID', 'PROCESSING', 'PROCESSED');
CREATE TYPE reward_types as ENUM ('%', 'pt');

CREATE TABLE accrual_orders (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_number TEXT UNIQUE NOT NULL,
    status accrual_order_status NOT NULL DEFAULT 'REGISTERED',
    accrual NUMERIC(12,2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE accrual_order_goods (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES accrual_orders(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    price NUMERIC(12,2) NOT NULL
);

CREATE TABLE reward_rules (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    match TEXT UNIQUE NOT NULL,
    reward NUMERIC(12,2) NOT NULL,
    reward_type reward_types NOT NULL, -- "%" или "pt"
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS reward_rules;
DROP TABLE IF EXISTS accrual_order_goods;
DROP TABLE IF EXISTS accrual_orders;
DROP TYPE IF EXISTS accrual_order_status;
