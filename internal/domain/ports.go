package domain

import "context"

// AIService defines the port for AI-related operations
type AIService interface {
	ParseTextToTransactions(ctx context.Context, text string) ([]Transaction, error)
}

// AuthService defines the port for authentication operations
type AuthService interface {
	ValidateToken(ctx context.Context, token string) (*AuthUser, error)
}

// TransactionRepository defines the port for transaction persistence
type TransactionRepository interface {
	SaveTransactions(ctx context.Context, transactions []Transaction) error
	GetTransactionByID(ctx context.Context, id int) (*Transaction, error)
	GetTransactions(ctx context.Context, limit, offset int) ([]Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *Transaction) error
	DeleteTransaction(ctx context.Context, id int) error
}

// TransactionService defines the port for transaction business logic
type TransactionService interface {
	SaveTransactions(ctx context.Context, transactions []Transaction) error
	GetTransactionByID(ctx context.Context, id int) (*Transaction, error)
	GetTransactions(ctx context.Context, limit, offset int) ([]Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *Transaction) error
	DeleteTransaction(ctx context.Context, id int) error
}
