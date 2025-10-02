-- Создание таблицы для учёта операций с баллами лояльности
CREATE TABLE IF NOT EXISTS loyalty_ledger (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_number VARCHAR(255),
    operation_type VARCHAR(20) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Создание индексов для таблицы операций
CREATE INDEX IF NOT EXISTS idx_loyalty_ledger_user_id ON loyalty_ledger(user_id);
CREATE INDEX IF NOT EXISTS idx_loyalty_ledger_operation_type ON loyalty_ledger(operation_type);
CREATE INDEX IF NOT EXISTS idx_loyalty_ledger_processed_at ON loyalty_ledger(processed_at);
CREATE INDEX IF NOT EXISTS idx_loyalty_ledger_order_number ON loyalty_ledger(order_number);

-- Ограничение на типы операций
ALTER TABLE loyalty_ledger ADD CONSTRAINT chk_operation_type 
CHECK (operation_type IN ('CREDIT', 'DEBIT'));

-- Ограничение на сумму (должна быть положительной)
ALTER TABLE loyalty_ledger ADD CONSTRAINT chk_amount_positive 
CHECK (amount > 0);
