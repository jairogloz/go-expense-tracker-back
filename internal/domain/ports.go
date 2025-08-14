package domain

import "context"

// AIService defines the port for AI-related operations
type AIService interface {
	ParseTextToTransactions(ctx context.Context, text string) ([]Transaction, error)
}

// TransactionRepository defines the port for transaction persistence
type TransactionRepository interface {
	SaveTransactions(ctx context.Context, transactions []Transaction) error
	GetTransactionByID(ctx context.Context, id int) (*Transaction, error)
	GetTransactions(ctx context.Context, limit, offset int) ([]Transaction, error)
}

// TransactionService defines the port for transaction business logic
type TransactionService interface {
	SaveTransactions(ctx context.Context, transactions []Transaction) error
	GetTransactionByID(ctx context.Context, id int) (*Transaction, error)
	GetTransactions(ctx context.Context, limit, offset int) ([]Transaction, error)
}
