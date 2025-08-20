package infra

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jairogloz/go-expense-tracker-back/internal/domain"
)

// PostgreSQLTransactionRepository implements the TransactionRepository interface
type PostgreSQLTransactionRepository struct {
	db *pgxpool.Pool
}

// NewPostgreSQLTransactionRepository creates a new PostgreSQL transaction repository
func NewPostgreSQLTransactionRepository(db *pgxpool.Pool) *PostgreSQLTransactionRepository {
	return &PostgreSQLTransactionRepository{
		db: db,
	}
}

// SaveTransactions saves multiple transactions to the database
func (r *PostgreSQLTransactionRepository) SaveTransactions(ctx context.Context, userID string, transactions []domain.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Prepare the insert statement
	stmt := `INSERT INTO transactions (user_id, account_id, amount, currency, category, type, date, description, sub_category) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	for _, transaction := range transactions {
		_, err := tx.Exec(ctx, stmt,
			userID,
			transaction.AccountID, // This can be nil
			transaction.Amount,
			transaction.Currency,
			transaction.Category,
			transaction.Type,
			transaction.Date,
			transaction.Description,
			transaction.SubCategory,
		)
		if err != nil {
			return fmt.Errorf("failed to insert transaction: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTransactionByID retrieves a transaction by its ID for a specific user
func (r *PostgreSQLTransactionRepository) GetTransactionByID(ctx context.Context, userID string, id int) (*domain.Transaction, error) {
	stmt := `SELECT id, user_id, account_id, amount, currency, category, type, date, description, sub_category 
			 FROM transactions WHERE id = $1 AND user_id = $2`

	var transaction domain.Transaction
	err := r.db.QueryRow(ctx, stmt, id, userID).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.AccountID,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.Category,
		&transaction.Type,
		&transaction.Date,
		&transaction.Description,
		&transaction.SubCategory,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Transaction not found or doesn't belong to user
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

// GetTransactions retrieves transactions for a specific user with pagination
func (r *PostgreSQLTransactionRepository) GetTransactions(ctx context.Context, userID string, limit, offset int) ([]domain.Transaction, error) {
	stmt := `SELECT id, user_id, account_id, amount, currency, category, type, date, description, sub_category 
			 FROM transactions WHERE user_id = $1 ORDER BY date DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, stmt, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.AccountID,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Category,
			&transaction.Type,
			&transaction.Date,
			&transaction.Description,
			&transaction.SubCategory,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return transactions, nil
}

// CreateTransactionsTable creates the transactions table if it doesn't exist
func (r *PostgreSQLTransactionRepository) CreateTransactionsTable(ctx context.Context) error {
	// Create the table with foreign key constraint from the start
	createTableStmt := `
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
		sub_category VARCHAR(100),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_transactions_account_id 
			FOREIGN KEY (account_id) REFERENCES accounts(id) 
			ON DELETE SET NULL
	);`

	_, err := r.db.Exec(ctx, createTableStmt)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}

	// Create indexes for better query performance
	indexStmt := `
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
	-- Sub-category index for filtering
	CREATE INDEX IF NOT EXISTS idx_transactions_sub_category ON transactions(sub_category);
	`

	_, err = r.db.Exec(ctx, indexStmt)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// UpdateTransaction updates an existing transaction for a specific user
func (r *PostgreSQLTransactionRepository) UpdateTransaction(ctx context.Context, userID string, transaction *domain.Transaction) error {
	stmt := `UPDATE transactions 
			 SET account_id = $3, amount = $4, currency = $5, category = $6, type = $7, date = $8, description = $9, sub_category = $10, updated_at = CURRENT_TIMESTAMP
			 WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(ctx, stmt,
		transaction.ID,
		userID,
		transaction.AccountID,
		transaction.Amount,
		transaction.Currency,
		transaction.Category,
		transaction.Type,
		transaction.Date,
		transaction.Description,
		transaction.SubCategory,
	)

	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("transaction with id %d not found or doesn't belong to user", transaction.ID)
	}

	return nil
}

// DeleteTransaction deletes a transaction by ID for a specific user
func (r *PostgreSQLTransactionRepository) DeleteTransaction(ctx context.Context, userID string, id int) error {
	stmt := `DELETE FROM transactions WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(ctx, stmt, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("transaction with id %d not found or doesn't belong to user", id)
	}

	return nil
}
