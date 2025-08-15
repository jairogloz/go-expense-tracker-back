package domain

import (
	"time"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

// Category represents predefined transaction categories
type Category string

const (
	// Expense categories
	CategoryFood          Category = "food"
	CategoryTransport     Category = "transport"
	CategoryUtilities     Category = "utilities"
	CategoryShopping      Category = "shopping"
	CategoryHealth        Category = "health"
	CategoryEducation     Category = "education"
	CategoryEntertainment Category = "entertainment"
	CategoryOther         Category = "other"

	// Income categories
	CategorySalary      Category = "salary"
	CategoryFreelance   Category = "freelance"
	CategoryInvestments Category = "investments"
	CategoryBonus       Category = "bonus"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID          int             `json:"id"`
	Amount      float64         `json:"amount"`
	Currency    string          `json:"currency"`
	Category    Category        `json:"category"`
	Type        TransactionType `json:"type"`
	Date        time.Time       `json:"date"`
	Description string          `json:"description,omitempty"`
}

// ParseInputRequest represents the request for parsing natural language input
type ParseInputRequest struct {
	Text string `json:"text" binding:"required"`
}

// ParseInputResponse represents the response after parsing input
type ParseInputResponse struct {
	Transactions []Transaction `json:"transactions"`
	Message      string        `json:"message,omitempty"`
}

// UpdateTransactionRequest represents the request for updating a transaction
type UpdateTransactionRequest struct {
	Amount      float64         `json:"amount" binding:"required,gt=0"`
	Currency    string          `json:"currency" binding:"required,len=3"`
	Category    Category        `json:"category" binding:"required"`
	Type        TransactionType `json:"type" binding:"required"`
	Date        time.Time       `json:"date" binding:"required"`
	Description string          `json:"description"`
}
