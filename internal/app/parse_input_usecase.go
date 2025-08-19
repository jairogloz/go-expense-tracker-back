package app

import (
	"context"
	"fmt"

	"github.com/jairogloz/go-expense-tracker-back/internal/domain"
)

// ParseInputUseCase handles the parsing of natural language input into transactions
type ParseInputUseCase struct {
	aiService          domain.AIService
	transactionService domain.TransactionService
	accountService     domain.AccountService
}

// NewParseInputUseCase creates a new parse input use case
func NewParseInputUseCase(aiService domain.AIService, transactionService domain.TransactionService, accountService domain.AccountService) *ParseInputUseCase {
	return &ParseInputUseCase{
		aiService:          aiService,
		transactionService: transactionService,
		accountService:     accountService,
	}
}

// Execute parses the input text and saves the resulting transactions
func (uc *ParseInputUseCase) Execute(ctx context.Context, request domain.ParseInputRequest) (*domain.ParseInputResponse, error) {
	// Extract user ID from context
	userID, ok := ctx.Value(domain.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user ID not found in context")
	}

	// Get user accounts for AI context
	// Todo: cache in the future so we don't hit the database every time
	accountsMap, err := uc.accountService.GetUserAccountsMap(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %w", err)
	}

	// Parse the text using AI service with account context
	transactions, err := uc.aiService.ParseTextToTransactionsWithAccounts(ctx, request.Text, accountsMap)
	if err != nil {
		return nil, err
	}

	// For transactions without account_id, try to assign default account
	defaultAccount, err := uc.accountService.GetDefaultAccount(ctx, userID)
	if err == nil && defaultAccount != nil {
		for i := range transactions {
			if transactions[i].AccountID == nil {
				transactions[i].AccountID = &defaultAccount.ID
			}
		}
	}

	// Save the transactions using transaction service
	if len(transactions) > 0 {
		if err := uc.transactionService.SaveTransactions(ctx, userID, transactions); err != nil {
			return nil, err
		}
	}

	response := &domain.ParseInputResponse{
		Transactions: transactions,
		Message:      "Successfully parsed and saved transactions",
	}

	return response, nil
}
