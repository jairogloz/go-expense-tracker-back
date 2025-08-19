package services

import (
	"context"
	"fmt"

	"github.com/jairogloz/go-expense-tracker-back/internal/app"
	"github.com/jairogloz/go-expense-tracker-back/internal/domain"
)

// AccountServiceImpl implements the AccountService interface
type AccountServiceImpl struct {
	repo           domain.AccountRepository
	balanceUseCase *app.AccountBalanceUseCase
}

// NewAccountService creates a new account service
func NewAccountService(repo domain.AccountRepository, balanceUseCase *app.AccountBalanceUseCase) *AccountServiceImpl {
	return &AccountServiceImpl{
		repo:           repo,
		balanceUseCase: balanceUseCase,
	}
}

// CreateAccount creates a new account for a user
func (s *AccountServiceImpl) CreateAccount(ctx context.Context, userID string, request domain.CreateAccountRequest) (*domain.Account, error) {
	// Validate account type
	if request.Type != domain.AccountTypeDebit && request.Type != domain.AccountTypeCredit {
		return nil, fmt.Errorf("invalid account type: %s", request.Type)
	}

	// Check if account name already exists for this user
	existing, err := s.repo.GetAccountByName(ctx, userID, request.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing account: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("account with name '%s' already exists", request.Name)
	}

	account := &domain.Account{
		UserID:         userID,
		Name:           request.Name,
		Type:           request.Type,
		InitialBalance: request.InitialBalance,
		IsDefault:      request.IsDefault,
	}

	if err := s.repo.CreateAccount(ctx, userID, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

// GetAccount retrieves an account by ID
func (s *AccountServiceImpl) GetAccount(ctx context.Context, userID string, accountID string) (*domain.Account, error) {
	account, err := s.repo.GetAccount(ctx, userID, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return nil, fmt.Errorf("account not found")
	}
	return account, nil
}

// GetDefaultAccount retrieves the default account for a user
func (s *AccountServiceImpl) GetDefaultAccount(ctx context.Context, userID string) (*domain.Account, error) {
	return s.repo.GetDefaultAccount(ctx, userID)
}

// GetUserAccounts retrieves all accounts for a user
func (s *AccountServiceImpl) GetUserAccounts(ctx context.Context, userID string) ([]domain.Account, error) {
	return s.repo.GetUserAccounts(ctx, userID)
}

// GetAccountBalance calculates and returns the current balance for an account
func (s *AccountServiceImpl) GetAccountBalance(ctx context.Context, userID string, accountID string) (*domain.AccountSummary, error) {
	return s.balanceUseCase.GetAccountBalance(ctx, userID, accountID)
}

// UpdateAccount updates an existing account
func (s *AccountServiceImpl) UpdateAccount(ctx context.Context, userID string, accountID string, request domain.UpdateAccountRequest) (*domain.Account, error) {
	// Get the existing account
	account, err := s.repo.GetAccount(ctx, userID, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return nil, fmt.Errorf("account not found")
	}

	// Update fields if provided
	if request.Name != nil {
		// Check if new name conflicts with existing accounts
		if *request.Name != account.Name {
			existing, err := s.repo.GetAccountByName(ctx, userID, *request.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to check existing account: %w", err)
			}
			if existing != nil {
				return nil, fmt.Errorf("account with name '%s' already exists", *request.Name)
			}
		}
		account.Name = *request.Name
	}

	if request.Type != nil {
		if *request.Type != domain.AccountTypeDebit && *request.Type != domain.AccountTypeCredit {
			return nil, fmt.Errorf("invalid account type: %s", *request.Type)
		}
		account.Type = *request.Type
	}

	if request.IsDefault != nil {
		account.IsDefault = *request.IsDefault
	}

	if err := s.repo.UpdateAccount(ctx, userID, account); err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return account, nil
}

// DeleteAccount deletes an account
func (s *AccountServiceImpl) DeleteAccount(ctx context.Context, userID string, accountID string) error {
	// Check if account exists
	account, err := s.repo.GetAccount(ctx, userID, accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return fmt.Errorf("account not found")
	}

	return s.repo.DeleteAccount(ctx, userID, accountID)
}

// SetDefaultAccount sets an account as the default for a user
func (s *AccountServiceImpl) SetDefaultAccount(ctx context.Context, userID string, accountID string) error {
	// Check if account exists
	account, err := s.repo.GetAccount(ctx, userID, accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return fmt.Errorf("account not found")
	}

	return s.repo.SetDefaultAccount(ctx, userID, accountID)
}

// GetUserAccountsMap returns a map of account names to IDs for AI processing
func (s *AccountServiceImpl) GetUserAccountsMap(ctx context.Context, userID string) (map[string]string, error) {
	accounts, err := s.repo.GetUserAccounts(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %w", err)
	}

	accountsMap := make(map[string]string)
	for _, account := range accounts {
		accountsMap[account.Name] = account.ID
	}

	return accountsMap, nil
}
