package services

import (
	"context"

	"github.com/jairogloz/go-expense-tracker-back/internal/domain"
)

// TransactionServiceImpl implements the TransactionService interface
type TransactionServiceImpl struct {
	repo domain.TransactionRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(repo domain.TransactionRepository) *TransactionServiceImpl {
	return &TransactionServiceImpl{
		repo: repo,
	}
}

// SaveTransactions saves multiple transactions for a specific user
func (s *TransactionServiceImpl) SaveTransactions(ctx context.Context, userID string, transactions []domain.Transaction) error {
	// Set the user ID for all transactions
	for i := range transactions {
		transactions[i].UserID = userID
	}
	return s.repo.SaveTransactions(ctx, userID, transactions)
}

// GetTransactionByID retrieves a transaction by its ID for a specific user
func (s *TransactionServiceImpl) GetTransactionByID(ctx context.Context, userID string, id int) (*domain.Transaction, error) {
	return s.repo.GetTransactionByID(ctx, userID, id)
}

// GetTransactions retrieves transactions for a specific user with pagination
func (s *TransactionServiceImpl) GetTransactions(ctx context.Context, userID string, limit, offset int) ([]domain.Transaction, error) {
	return s.repo.GetTransactions(ctx, userID, limit, offset)
}

// UpdateTransaction updates an existing transaction for a specific user
func (s *TransactionServiceImpl) UpdateTransaction(ctx context.Context, userID string, transaction *domain.Transaction) error {
	// Ensure the transaction belongs to the user
	transaction.UserID = userID
	return s.repo.UpdateTransaction(ctx, userID, transaction)
}

// DeleteTransaction deletes a transaction by ID for a specific user
func (s *TransactionServiceImpl) DeleteTransaction(ctx context.Context, userID string, id int) error {
	return s.repo.DeleteTransaction(ctx, userID, id)
}
