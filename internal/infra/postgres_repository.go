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
func (r *PostgreSQLTransactionRepository) SaveTransactions(ctx context.Context, transactions []domain.Transaction) error {
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
	stmt := `INSERT INTO transactions (amount, currency, category, type, date, vendor, description) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, transaction := range transactions {
		_, err := tx.Exec(ctx, stmt,
			transaction.Amount,
			transaction.Currency,
			transaction.Category,
			transaction.Type,
			transaction.Date,
			transaction.Vendor,
			transaction.Description,
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

// GetTransactionByID retrieves a transaction by its ID
func (r *PostgreSQLTransactionRepository) GetTransactionByID(ctx context.Context, id int) (*domain.Transaction, error) {
	stmt := `SELECT id, amount, currency, category, type, date, vendor, description 
			 FROM transactions WHERE id = $1`

	var transaction domain.Transaction
	err := r.db.QueryRow(ctx, stmt, id).Scan(
		&transaction.ID,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.Category,
		&transaction.Type,
		&transaction.Date,
		&transaction.Vendor,
		&transaction.Description,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Transaction not found
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

// GetTransactions retrieves transactions with pagination
func (r *PostgreSQLTransactionRepository) GetTransactions(ctx context.Context, limit, offset int) ([]domain.Transaction, error) {
	stmt := `SELECT id, amount, currency, category, type, date, vendor, description 
			 FROM transactions ORDER BY date DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Category,
			&transaction.Type,
			&transaction.Date,
			&transaction.Vendor,
			&transaction.Description,
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
	stmt := `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		amount DECIMAL(12,2) NOT NULL,
		currency VARCHAR(3) NOT NULL DEFAULT 'USD',
		category VARCHAR(50) NOT NULL,
		type VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
		date TIMESTAMP NOT NULL,
		vendor VARCHAR(255) NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	-- Create index on date for better query performance
	CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
	-- Create index on type for filtering
	CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
	-- Create index on category for filtering
	CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category);
	`

	_, err := r.db.Exec(ctx, stmt)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}

	return nil
}
