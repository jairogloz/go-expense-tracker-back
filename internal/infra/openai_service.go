package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jairogloz/go-expense-tracker-back/internal/domain"
	"github.com/sashabaranov/go-openai"
)

// OpenAIService implements the AIService interface
type OpenAIService struct {
	client *openai.Client
}

// NewOpenAIService creates a new OpenAI service
func NewOpenAIService(apiKey string) *OpenAIService {
	client := openai.NewClient(apiKey)
	return &OpenAIService{
		client: client,
	}
}

// ParseTextToTransactions parses natural language text into structured transactions
func (s *OpenAIService) ParseTextToTransactions(ctx context.Context, text string) ([]domain.Transaction, error) {
	systemPrompt := `You are a financial transaction parser. Parse the given text into structured transaction data.

Available categories:
- Expense: food, transport, utilities, shopping, health, education, entertainment, other
- Income: salary, freelance, investments, bonus

Return a JSON array of transactions with the following structure:
{
  "transactions": [
    {
      "amount": 25.50,
      "currency": "MXN",
      "category": "food",
      "type": "expense",
      "date": "2024-01-15T12:00:00Z",
      "vendor": "Restaurant Name",
      "description": "Lunch at restaurant"
    }
  ]
}

Rules:
1. If no date is specified, use the current date
2. Default currency is MXN if not specified
3. Amount should be positive (the type field indicates income/expense)
4. Choose the most appropriate category from the available list
5. Extract vendor name from the text
6. If multiple transactions are mentioned, create separate objects for each

Parse this text:`

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			},
		},
		MaxTokens:   1000,
		Temperature: 0.1,
	}

	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI API")
	}

	content := resp.Choices[0].Message.Content

	// Parse the JSON response
	var response struct {
		Transactions []struct {
			Amount      float64 `json:"amount"`
			Currency    string  `json:"currency"`
			Category    string  `json:"category"`
			Type        string  `json:"type"`
			Date        string  `json:"date"`
			Vendor      string  `json:"vendor"`
			Description string  `json:"description"`
		} `json:"transactions"`
	}

	if err := json.Unmarshal([]byte(content), &response); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w, content: %s", err, content)
	}

	// Convert to domain transactions
	var transactions []domain.Transaction
	for _, t := range response.Transactions {
		// Parse date
		date, err := time.Parse(time.RFC3339, t.Date)
		if err != nil {
			// If parsing fails, use current time
			date = time.Now()
		}

		// Validate and convert type
		var transactionType domain.TransactionType
		switch t.Type {
		case "income":
			transactionType = domain.Income
		case "expense":
			transactionType = domain.Expense
		default:
			transactionType = domain.Expense // default to expense
		}

		// Validate and convert category
		category := domain.Category(t.Category)
		// You could add validation here to ensure category is valid

		transaction := domain.Transaction{
			Amount:      t.Amount,
			Currency:    t.Currency,
			Category:    category,
			Type:        transactionType,
			Date:        date,
			Vendor:      t.Vendor,
			Description: t.Description,
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
