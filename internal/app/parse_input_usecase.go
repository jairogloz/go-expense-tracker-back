package app

import (
	"context"

	"github.com/jairogloz/go-expense-tracker-back/internal/domain"
)

// ParseInputUseCase handles the parsing of natural language input into transactions
type ParseInputUseCase struct {
	aiService          domain.AIService
	transactionService domain.TransactionService
}

// NewParseInputUseCase creates a new parse input use case
func NewParseInputUseCase(aiService domain.AIService, transactionService domain.TransactionService) *ParseInputUseCase {
	return &ParseInputUseCase{
		aiService:          aiService,
		transactionService: transactionService,
	}
}

// Execute parses the input text and saves the resulting transactions
func (uc *ParseInputUseCase) Execute(ctx context.Context, request domain.ParseInputRequest) (*domain.ParseInputResponse, error) {
	// Parse the text using AI service
	transactions, err := uc.aiService.ParseTextToTransactions(ctx, request.Text)
	if err != nil {
		return nil, err
	}

	// Save the transactions using transaction service
	if len(transactions) > 0 {
		if err := uc.transactionService.SaveTransactions(ctx, transactions); err != nil {
			return nil, err
		}
	}

	response := &domain.ParseInputResponse{
		Transactions: transactions,
		Message:      "Successfully parsed and saved transactions",
	}

	return response, nil
}
