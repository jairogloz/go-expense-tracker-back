package domain

import (
	"context"
	"time"
)

// AIService defines the port for AI-related operations
type AIService interface {
	ParseTextToTransactions(ctx context.Context, text string) ([]Transaction, error)
	ParseTextToTransactionsWithAccounts(ctx context.Context, text string, accountsMap map[string]string) ([]Transaction, error)
}

// AuthService defines the port for authentication operations
type AuthService interface {
	ValidateToken(ctx context.Context, token string) (*AuthUser, error)
}

// AccountRepository defines the port for account persistence
type AccountRepository interface {
	CreateAccount(ctx context.Context, userID string, account *Account) error
	GetAccount(ctx context.Context, userID string, accountID string) (*Account, error)
	GetAccountByName(ctx context.Context, userID string, name string) (*Account, error)
	GetUserAccounts(ctx context.Context, userID string) ([]Account, error)
	GetDefaultAccount(ctx context.Context, userID string) (*Account, error)
	UpdateAccount(ctx context.Context, userID string, account *Account) error
	DeleteAccount(ctx context.Context, userID string, accountID string) error
	SetDefaultAccount(ctx context.Context, userID string, accountID string) error
	SumTransactionsSince(ctx context.Context, accountID string, since time.Time) (float64, error)
}

// AccountService defines the port for account business logic
type AccountService interface {
	CreateAccount(ctx context.Context, userID string, request CreateAccountRequest) (*Account, error)
	GetAccount(ctx context.Context, userID string, accountID string) (*Account, error)
	GetDefaultAccount(ctx context.Context, userID string) (*Account, error)
	GetUserAccounts(ctx context.Context, userID string) ([]Account, error)
	GetAccountBalance(ctx context.Context, userID string, accountID string) (*AccountSummary, error)
	UpdateAccount(ctx context.Context, userID string, accountID string, request UpdateAccountRequest) (*Account, error)
	DeleteAccount(ctx context.Context, userID string, accountID string) error
	SetDefaultAccount(ctx context.Context, userID string, accountID string) error
	GetUserAccountsMap(ctx context.Context, userID string) (map[string]string, error) // name -> id mapping
}

// TransactionRepository defines the port for transaction persistence
type TransactionRepository interface {
	SaveTransactions(ctx context.Context, userID string, transactions []Transaction) error
	GetTransactionByID(ctx context.Context, userID string, id int) (*Transaction, error)
	GetTransactions(ctx context.Context, userID string, limit, offset int) ([]Transaction, error)
	UpdateTransaction(ctx context.Context, userID string, transaction *Transaction) error
	DeleteTransaction(ctx context.Context, userID string, id int) error
}

// TransactionService defines the port for transaction business logic
type TransactionService interface {
	SaveTransactions(ctx context.Context, userID string, transactions []Transaction) error
	GetTransactionByID(ctx context.Context, userID string, id int) (*Transaction, error)
	GetTransactions(ctx context.Context, userID string, limit, offset int) ([]Transaction, error)
	UpdateTransaction(ctx context.Context, userID string, transaction *Transaction) error
	DeleteTransaction(ctx context.Context, userID string, id int) error
}
