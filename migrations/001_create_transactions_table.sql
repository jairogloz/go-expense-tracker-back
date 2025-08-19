-- Migration: 001_create_transactions_table.sql
-- Description: Create the transactions table with proper indexes and constraints

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    account_id UUID,
    amount DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    category VARCHAR(50) NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    date TIMESTAMP NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
-- Primary index for user-scoped queries (most important)
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
-- Composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_transactions_user_date ON transactions(user_id, date DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_user_category ON transactions(user_id, category);
CREATE INDEX IF NOT EXISTS idx_transactions_user_type ON transactions(user_id, type);
CREATE INDEX IF NOT EXISTS idx_transactions_user_created_at ON transactions(user_id, created_at DESC);
-- Account-specific indexes for balance calculations
CREATE INDEX IF NOT EXISTS idx_transactions_account_date ON transactions(account_id, date);
CREATE INDEX IF NOT EXISTS idx_transactions_account_amount ON transactions(account_id, amount);
-- Individual indexes for filtering
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_transactions_updated_at 
    BEFORE UPDATE ON transactions 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE transactions IS 'Stores financial transactions parsed from natural language input';
COMMENT ON COLUMN transactions.amount IS 'Transaction amount (always positive, type indicates income/expense)';
COMMENT ON COLUMN transactions.currency IS 'ISO 4217 currency code (e.g., USD, EUR)';
COMMENT ON COLUMN transactions.category IS 'Predefined category for the transaction';
COMMENT ON COLUMN transactions.type IS 'Transaction type: income or expense';
COMMENT ON COLUMN transactions.date IS 'Date when the transaction occurred';
COMMENT ON COLUMN transactions.description IS 'Additional details about the transaction';
