-- Migration: 002_create_accounts_table.sql
-- Description: Create the accounts table with proper indexes and constraints

-- Create accounts table
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('debit', 'credit')),
    initial_balance DECIMAL(12,2) NOT NULL DEFAULT 0,
    last_calculated_balance DECIMAL(12,2) DEFAULT 0,
    last_calculated_at TIMESTAMP,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name) -- Prevent duplicate account names per user
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_accounts_user_default ON accounts(user_id, is_default);
CREATE INDEX IF NOT EXISTS idx_accounts_user_name ON accounts(user_id, name);

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_accounts_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_accounts_updated_at BEFORE UPDATE ON accounts
    FOR EACH ROW EXECUTE FUNCTION update_accounts_updated_at_column();

-- Add foreign key constraint to transactions table
ALTER TABLE transactions 
ADD CONSTRAINT fk_transactions_account 
FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE SET NULL;
