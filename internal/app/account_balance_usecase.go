package app

import (
	"context"
	"fmt"

	"github.com/jairogloz/go-expense-tracker-back/internal/domain"
)

// AccountBalanceUseCase handles account balance calculations
type AccountBalanceUseCase struct {
	accountRepo domain.AccountRepository
}

// NewAccountBalanceUseCase creates a new account balance use case
func NewAccountBalanceUseCase(accountRepo domain.AccountRepository, transactionRepo domain.TransactionRepository) *AccountBalanceUseCase {
	return &AccountBalanceUseCase{
		accountRepo: accountRepo,
	}
}

// GetAccountBalance calculates the current balance for an account
func (uc *AccountBalanceUseCase) GetAccountBalance(ctx context.Context, userID string, accountID string) (*domain.AccountSummary, error) {
	// Get the account
	account, err := uc.accountRepo.GetAccount(ctx, userID, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Calculate balance based on initial balance + transaction sum
	transactionSum, err := uc.accountRepo.SumTransactionsSince(ctx, accountID, account.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to sum transactions: %w", err)
	}

	// For credit accounts, expenses increase debt (negative balance gets more negative)
	// For debit accounts, expenses decrease balance
	currentBalance := account.InitialBalance
	if account.Type == domain.AccountTypeCredit {
		// Credit card logic: expenses subtract from balance (increase debt), payments add to balance (decrease debt)
		currentBalance += transactionSum
	} else {
		// Debit account logic: expenses subtract, income adds
		currentBalance += transactionSum
	}

	return &domain.AccountSummary{
		Account:        account,
		CurrentBalance: currentBalance,
	}, nil
}
