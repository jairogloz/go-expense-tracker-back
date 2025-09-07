package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jairogloz/go-expense-tracker-back/internal/domain"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// OpenAIService implements the AIService interface
type OpenAIService struct {
	client *openai.Client
	logger *zap.Logger
}

// NewOpenAIService creates a new OpenAI service
func NewOpenAIService(apiKey string, logger *zap.Logger) *OpenAIService {
	client := openai.NewClient(apiKey)
	return &OpenAIService{
		client: client,
		logger: logger,
	}
}

// ParseTextToTransactions parses natural language text into structured transactions
func (s *OpenAIService) ParseTextToTransactions(ctx context.Context, text string) ([]domain.Transaction, error) {
	// Extract context info for logging
	userID := ctx.Value(domain.UserIDKey)
	requestID := ctx.Value(domain.ContextKey("request_id"))

	s.logger.Info("Starting OpenAI transaction parsing",
		zap.Any("user_id", userID),
		zap.Any("request_id", requestID),
		zap.Int("input_text_length", len(text)),
		zap.Int("estimated_transactions", strings.Count(text, ",")+1), // Rough estimate
	)

	now := time.Now().UTC().Format(time.RFC3339)
	systemPrompt := fmt.Sprintf(`You are a financial transaction parser. Parse the given text into structured transaction data.

Available categories:
- Expense: guilt_free_spending, fixed_costs, investments, savings_goals
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
      "description": "Electricity bill"
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

	// Calculate appropriate token limit based on input
	maxTokens := s.calculateMaxTokens(text)

	s.logger.Debug("Calculated token requirements",
		zap.Any("user_id", userID),
		zap.Any("request_id", requestID),
		zap.Int("calculated_max_tokens", maxTokens),
		zap.Int("estimated_transactions", strings.Count(text, ",")+1),
	)

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
		MaxTokens:   maxTokens,
		Temperature: 0.1,
	}

	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		s.logger.Error("OpenAI API call failed",
			zap.Any("user_id", userID),
			zap.Any("request_id", requestID),
			zap.Error(err),
			zap.String("model", string(req.Model)),
			zap.Int("max_tokens", req.MaxTokens),
			zap.Float32("temperature", req.Temperature),
			zap.Int("input_text_length", len(text)),
		)
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(resp.Choices) == 0 {
		s.logger.Error("OpenAI API returned no response choices",
			zap.Any("user_id", userID),
			zap.Any("request_id", requestID),
			zap.String("model", string(req.Model)),
		)
		return nil, fmt.Errorf("no response from OpenAI API")
	}

	content := resp.Choices[0].Message.Content

	s.logger.Debug("OpenAI API response received",
		zap.Any("user_id", userID),
		zap.Any("request_id", requestID),
		zap.Int("response_length", len(content)),
		zap.String("finish_reason", string(resp.Choices[0].FinishReason)),
		zap.Int("prompt_tokens", resp.Usage.PromptTokens),
		zap.Int("completion_tokens", resp.Usage.CompletionTokens),
		zap.Int("total_tokens", resp.Usage.TotalTokens),
		zap.String("raw_response", content), // Only in debug level due to potentially sensitive data
	)

	// Check if response was truncated due to token limit
	if resp.Choices[0].FinishReason == openai.FinishReasonLength {
		s.logger.Warn("OpenAI response was truncated due to token limit",
			zap.Any("user_id", userID),
			zap.Any("request_id", requestID),
			zap.Int("max_tokens_requested", maxTokens),
			zap.Int("completion_tokens_used", resp.Usage.CompletionTokens),
			zap.Int("response_length", len(content)),
		)
	}

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
		} `json:"transactions"`
	}

	if err := json.Unmarshal([]byte(content), &response); err != nil {
		// Check if this is due to truncated response
		isTruncated := resp.Choices[0].FinishReason == openai.FinishReasonLength
		s.logger.Error("Failed to parse OpenAI JSON response",
			zap.Any("user_id", userID),
			zap.Any("request_id", requestID),
			zap.Error(err),
			zap.Bool("response_truncated", isTruncated),
			zap.String("finish_reason", string(resp.Choices[0].FinishReason)),
			zap.String("raw_response", content),
			zap.Int("response_length", len(content)),
		)

		if isTruncated {
			return nil, fmt.Errorf("OpenAI response was truncated due to token limit. Consider reducing input size or increasing MaxTokens. Original error: %w", err)
		}

		return nil, fmt.Errorf("failed to parse OpenAI response: %w, content: %s", err, content)
	}

	s.logger.Info("Successfully parsed OpenAI response",
		zap.Any("user_id", userID),
		zap.Any("request_id", requestID),
		zap.Int("transactions_count", len(response.Transactions)),
	)

	// Convert to domain transactions
	var transactions []domain.Transaction
	for i, t := range response.Transactions {
		// Parse date
		date, err := time.Parse(time.RFC3339, t.Date)
		if err != nil {
			s.logger.Warn("Failed to parse transaction date, using current time",
				zap.Any("user_id", userID),
				zap.Any("request_id", requestID),
				zap.Int("transaction_index", i),
				zap.String("original_date", t.Date),
				zap.Error(err),
			)
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
		}

		transactions = append(transactions, transaction)
	}

	s.logger.Info("Successfully converted transactions",
		zap.Any("user_id", userID),
		zap.Any("request_id", requestID),
		zap.Int("final_transactions_count", len(transactions)),
	)

	return transactions, nil
}

// calculateMaxTokens estimates the required tokens based on input characteristics
func (s *OpenAIService) calculateMaxTokens(inputText string) int {
	// Rough estimation:
	// - System prompt: ~1000 tokens
	// - Input text: ~1 token per 4 characters
	// - Each transaction in response: ~150-200 tokens

	estimatedTransactions := strings.Count(inputText, ",") + 1 // Rough estimate based on commas
	if estimatedTransactions < 1 {
		estimatedTransactions = 1
	}

	baseTokens := 1000                            // System prompt
	inputTokens := len(inputText) / 4             // Input text estimation
	responseTokens := estimatedTransactions * 200 // Response estimation (conservative)

	totalEstimate := baseTokens + inputTokens + responseTokens

	// Add 20% buffer and ensure minimum/maximum bounds
	withBuffer := int(float64(totalEstimate) * 1.2)

	// Bounds: minimum 1500, maximum 4000 (to stay within model limits)
	if withBuffer < 1500 {
		withBuffer = 1500
	}
	if withBuffer > 4000 {
		withBuffer = 4000
	}

	return withBuffer
}
