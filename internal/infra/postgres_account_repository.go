package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jairogloz/go-expense-tracker-back/internal/domain"
)

// PostgreSQLAccountRepository implements the AccountRepository interface
type PostgreSQLAccountRepository struct {
	db *pgxpool.Pool
}

// NewPostgreSQLAccountRepository creates a new PostgreSQL account repository
func NewPostgreSQLAccountRepository(db *pgxpool.Pool) *PostgreSQLAccountRepository {
	return &PostgreSQLAccountRepository{
		db: db,
	}
}

// CreateAccount creates a new account
func (r *PostgreSQLAccountRepository) CreateAccount(ctx context.Context, userID string, account *domain.Account) error {
	// If this is the first account for the user, make it default
	var userAccountCount int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM accounts WHERE user_id = $1", userID).Scan(&userAccountCount)
	if err != nil {
		return fmt.Errorf("failed to count user accounts: %w", err)
	}

	isDefault := account.IsDefault || userAccountCount == 0

	// If setting as default, unset other defaults first
	if isDefault {
		_, err = r.db.Exec(ctx, "UPDATE accounts SET is_default = FALSE WHERE user_id = $1", userID)
		if err != nil {
			return fmt.Errorf("failed to unset other default accounts: %w", err)
		}
	}

	stmt := `INSERT INTO accounts (user_id, name, type, initial_balance, last_calculated_balance, is_default) 
			 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`

	err = r.db.QueryRow(ctx, stmt,
		userID,
		account.Name,
		account.Type,
		account.InitialBalance,
		account.InitialBalance, // Set initial as last calculated
		isDefault,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	account.UserID = userID
	account.IsDefault = isDefault
	account.LastCalculatedBalance = account.InitialBalance
	account.LastCalculatedAt = account.CreatedAt

	return nil
}

// GetAccount retrieves an account by ID for a specific user
func (r *PostgreSQLAccountRepository) GetAccount(ctx context.Context, userID string, accountID string) (*domain.Account, error) {
	stmt := `SELECT id, user_id, name, type, initial_balance, last_calculated_balance, 
			 last_calculated_at, is_default, created_at, updated_at 
			 FROM accounts WHERE id = $1 AND user_id = $2`

	var account domain.Account
	var lastCalculatedAt *time.Time

	err := r.db.QueryRow(ctx, stmt, accountID, userID).Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&account.Type,
		&account.InitialBalance,
		&account.LastCalculatedBalance,
		&lastCalculatedAt,
		&account.IsDefault,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Account not found or doesn't belong to user
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	if lastCalculatedAt != nil {
		account.LastCalculatedAt = *lastCalculatedAt
	}

	return &account, nil
}

// GetAccountByName retrieves an account by name for a specific user
func (r *PostgreSQLAccountRepository) GetAccountByName(ctx context.Context, userID string, name string) (*domain.Account, error) {
	stmt := `SELECT id, user_id, name, type, initial_balance, last_calculated_balance, 
			 last_calculated_at, is_default, created_at, updated_at 
			 FROM accounts WHERE name = $1 AND user_id = $2`

	var account domain.Account
	var lastCalculatedAt *time.Time

	err := r.db.QueryRow(ctx, stmt, name, userID).Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&account.Type,
		&account.InitialBalance,
		&account.LastCalculatedBalance,
		&lastCalculatedAt,
		&account.IsDefault,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Account not found
		}
		return nil, fmt.Errorf("failed to get account by name: %w", err)
	}

	if lastCalculatedAt != nil {
		account.LastCalculatedAt = *lastCalculatedAt
	}

	return &account, nil
}

// GetUserAccounts retrieves all accounts for a user
func (r *PostgreSQLAccountRepository) GetUserAccounts(ctx context.Context, userID string) ([]domain.Account, error) {
	stmt := `SELECT id, user_id, name, type, initial_balance, last_calculated_balance, 
			 last_calculated_at, is_default, created_at, updated_at 
			 FROM accounts WHERE user_id = $1 ORDER BY is_default DESC, created_at ASC`

	rows, err := r.db.Query(ctx, stmt, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var account domain.Account
		var lastCalculatedAt *time.Time

		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Name,
			&account.Type,
			&account.InitialBalance,
			&account.LastCalculatedBalance,
			&lastCalculatedAt,
			&account.IsDefault,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}

		if lastCalculatedAt != nil {
			account.LastCalculatedAt = *lastCalculatedAt
		}

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return accounts, nil
}

// GetDefaultAccount retrieves the default account for a user
func (r *PostgreSQLAccountRepository) GetDefaultAccount(ctx context.Context, userID string) (*domain.Account, error) {
	stmt := `SELECT id, user_id, name, type, initial_balance, last_calculated_balance, 
			 last_calculated_at, is_default, created_at, updated_at 
			 FROM accounts WHERE user_id = $1 AND is_default = TRUE`

	var account domain.Account
	var lastCalculatedAt *time.Time

	err := r.db.QueryRow(ctx, stmt, userID).Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&account.Type,
		&account.InitialBalance,
		&account.LastCalculatedBalance,
		&lastCalculatedAt,
		&account.IsDefault,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No default account found
		}
		return nil, fmt.Errorf("failed to get default account: %w", err)
	}

	if lastCalculatedAt != nil {
		account.LastCalculatedAt = *lastCalculatedAt
	}

	return &account, nil
}

// UpdateAccount updates an existing account
func (r *PostgreSQLAccountRepository) UpdateAccount(ctx context.Context, userID string, account *domain.Account) error {
	// If setting as default, unset other defaults first
	if account.IsDefault {
		_, err := r.db.Exec(ctx, "UPDATE accounts SET is_default = FALSE WHERE user_id = $1 AND id != $2", userID, account.ID)
		if err != nil {
			return fmt.Errorf("failed to unset other default accounts: %w", err)
		}
	}

	stmt := `UPDATE accounts 
			 SET name = $3, type = $4, is_default = $5, updated_at = CURRENT_TIMESTAMP
			 WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(ctx, stmt,
		account.ID,
		userID,
		account.Name,
		account.Type,
		account.IsDefault,
	)

	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("account with id %s not found or doesn't belong to user", account.ID)
	}

	return nil
}

// DeleteAccount deletes an account by ID
func (r *PostgreSQLAccountRepository) DeleteAccount(ctx context.Context, userID string, accountID string) error {
	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Set account_id to NULL for all transactions of this account
	_, err = tx.Exec(ctx, "UPDATE transactions SET account_id = NULL WHERE account_id = $1", accountID)
	if err != nil {
		return fmt.Errorf("failed to unlink transactions: %w", err)
	}

	// Delete the account
	result, err := tx.Exec(ctx, "DELETE FROM accounts WHERE id = $1 AND user_id = $2", accountID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("account with id %s not found or doesn't belong to user", accountID)
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// SetDefaultAccount sets an account as the default for a user
func (r *PostgreSQLAccountRepository) SetDefaultAccount(ctx context.Context, userID string, accountID string) error {
	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Unset all other defaults
	_, err = tx.Exec(ctx, "UPDATE accounts SET is_default = FALSE WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("failed to unset other default accounts: %w", err)
	}

	// Set the new default
	result, err := tx.Exec(ctx, "UPDATE accounts SET is_default = TRUE WHERE id = $1 AND user_id = $2", accountID, userID)
	if err != nil {
		return fmt.Errorf("failed to set default account: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("account with id %s not found or doesn't belong to user", accountID)
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// SumTransactionsSince calculates the sum of transactions for an account since a given date
func (r *PostgreSQLAccountRepository) SumTransactionsSince(ctx context.Context, accountID string, since time.Time) (float64, error) {
	stmt := `SELECT 
			 COALESCE(
				 SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END), 
				 0
			 ) 
			 FROM transactions 
			 WHERE account_id = $1 AND created_at >= $2`

	var sum float64
	err := r.db.QueryRow(ctx, stmt, accountID, since).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("failed to sum transactions: %w", err)
	}

	return sum, nil
}

// CreateAccountsTable creates the accounts table if it doesn't exist
func (r *PostgreSQLAccountRepository) CreateAccountsTable(ctx context.Context) error {
	stmt := `
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
		UNIQUE(user_id, name)
	);
	
	-- Create indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
	CREATE INDEX IF NOT EXISTS idx_accounts_user_default ON accounts(user_id, is_default);
	CREATE INDEX IF NOT EXISTS idx_accounts_user_name ON accounts(user_id, name);
	
	-- Create trigger function if it doesn't exist
	CREATE OR REPLACE FUNCTION update_accounts_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ language 'plpgsql';
	
	-- Create trigger if it doesn't exist
	DROP TRIGGER IF EXISTS update_accounts_updated_at ON accounts;
	CREATE TRIGGER update_accounts_updated_at BEFORE UPDATE ON accounts
		FOR EACH ROW EXECUTE FUNCTION update_accounts_updated_at_column();
	`

	_, err := r.db.Exec(ctx, stmt)
	if err != nil {
		return fmt.Errorf("failed to create accounts table: %w", err)
	}

	return nil
}
