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

// SaveTransactions saves multiple transactions
func (s *TransactionServiceImpl) SaveTransactions(ctx context.Context, transactions []domain.Transaction) error {
	return s.repo.SaveTransactions(ctx, transactions)
}

// GetTransactionByID retrieves a transaction by its ID
func (s *TransactionServiceImpl) GetTransactionByID(ctx context.Context, id int) (*domain.Transaction, error) {
	return s.repo.GetTransactionByID(ctx, id)
}

// GetTransactions retrieves transactions with pagination
func (s *TransactionServiceImpl) GetTransactions(ctx context.Context, limit, offset int) ([]domain.Transaction, error) {
	return s.repo.GetTransactions(ctx, limit, offset)
}
