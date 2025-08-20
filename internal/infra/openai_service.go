package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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
	now := time.Now().UTC().Format(time.RFC3339)
	systemPrompt := fmt.Sprintf(`You are a financial transaction parser. Parse the given text into structured transaction data.

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
      "description": "Lunch at restaurant"
    }
  ]
}

Rules:
1. If no date and time is specified, use the current date and time.
2. Default currency is MXN if not specified
3. Amount should be positive (the type field indicates income/expense)
4. Choose the most appropriate category from the available list
5. If multiple transactions are mentioned, create separate objects for each
6. ALWAYS preserve the original language in descriptions - DO NOT TRANSLATE
7. Do NOT include the amount or currency in the description field - keep descriptions focused on what the transaction was for
8. Consider this date as the current date: %s

Parse this text:`, now)

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
			Description: t.Description,
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// ParseTextToTransactionsWithAccounts parses natural language text into structured transactions with account context
func (s *OpenAIService) ParseTextToTransactionsWithAccounts(ctx context.Context, text string, accountsMap map[string]string) ([]domain.Transaction, error) {
	now := time.Now().UTC().Format(time.RFC3339)

	// Build account information for the prompt
	accountInfo := "Available accounts:\n"
	defaultAccountID := ""
	for name, id := range accountsMap {
		accountInfo += fmt.Sprintf("- %s (ID: %s)\n", name, id)
		if defaultAccountID == "" {
			defaultAccountID = id // Use first account as fallback default
		}
	}

	systemPrompt := fmt.Sprintf(`You are a financial transaction parser. Parse the given text into structured transaction data.

Available categories:
- Expense: guilt_free, fixed_costs, investments and savings
- Income: salary, freelance, investments, bonus

Subcategories for each category are:
{
  "fixed_costs": [
    "house_payments",
    "utilities",
    "internet_and_phone",
    "insurance",
    "loan_payments",
    "subscriptions"
  ],
  "investments": [
    "afore",
    "cetes",
    "mutual_funds",
    "retirement_savings",
    "stocks",
    "real_estate",
    "crypto"
  ],
  "savings_goals": [
    "emergency_fund",
    "vacation_savings"
  ],
  "guilt_free_spending": [
    "dining_out",
    "coffee_and_snacks",
    "clothing_and_accessories",
    "hobbies",
    "entertainment",
    "fitness_and_wellness",
    "gifts_and_celebrations",
    "travel",
    "gadgets_and_tech"
  ]
}

%s

Return a JSON array of transactions with the following structure:
{
  "transactions": [
    {
      "amount": 25.50,
      "currency": "MXN",
      "category": "fixed_costs",
			"sub_category": "utilities",
      "type": "expense",
      "date": "2024-01-15T12:00:00Z",
      "description": "Electricity bill",
      "account_name": "Credit Card"
    }
  ]
}

Rules:
1. If no date and time is specified, use the current date and time.
2. Default currency is MXN if not specified
3. Amount should be positive (the type field indicates income/expense)
4. Choose the most appropriate category from the available list
5. If multiple transactions are mentioned, create separate objects for each
6. ALWAYS preserve the original language in descriptions - DO NOT TRANSLATE
7. For account_name: try to match account names mentioned in the text to the available accounts. If no account is mentioned or can be determined, use "default"
8. Do NOT include the amount or currency in the description field - keep descriptions focused on what the transaction was for
9. Consider this date as the current date: %s

Parse this text:`, accountInfo, now)

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
			SubCategory string  `json:"sub_category"`
			Type        string  `json:"type"`
			Date        string  `json:"date"`
			Description string  `json:"description"`
			AccountName string  `json:"account_name"`
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

		// Map account name to account ID
		var accountID *string
		if t.AccountName != "" && t.AccountName != "default" {
			// Try to find matching account ID by name
			if id, exists := accountsMap[t.AccountName]; exists {
				accountID = &id
			} else {
				// If account name doesn't match exactly, try case-insensitive match
				for name, id := range accountsMap {
					if strings.EqualFold(name, t.AccountName) {
						accountID = &id
						break
					}
				}
			}
		}

		// If no account was matched or "default" was specified, use default account
		if accountID == nil && defaultAccountID != "" {
			accountID = &defaultAccountID
		}

		// Handle SubCategory - only set if not empty
		var subCategory *string
		if t.SubCategory != "" {
			subCategory = &t.SubCategory
		}

		transaction := domain.Transaction{
			Amount:      t.Amount,
			Currency:    t.Currency,
			Category:    category,
			SubCategory: subCategory,
			Type:        transactionType,
			Date:        date,
			Description: t.Description,
			AccountID:   accountID,
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
