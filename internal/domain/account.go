package domain

import (
	"time"
)

// AccountType represents the type of account
type AccountType string

const (
	AccountTypeDebit  AccountType = "debit"
	AccountTypeCredit AccountType = "credit"
)

// Account represents a financial account
type Account struct {
	ID                    string      `json:"id"`
	UserID                string      `json:"user_id"`
	Name                  string      `json:"name"`
	Type                  AccountType `json:"type"`
	InitialBalance        float64     `json:"initial_balance"`
	LastCalculatedBalance float64     `json:"last_calculated_balance"`
	LastCalculatedAt      time.Time   `json:"last_calculated_at"`
	IsDefault             bool        `json:"is_default"`
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
}

// CreateAccountRequest represents the request for creating an account
type CreateAccountRequest struct {
	Name           string      `json:"name" binding:"required"`
	Type           AccountType `json:"type" binding:"required"`
	InitialBalance float64     `json:"initial_balance"`
	IsDefault      bool        `json:"is_default"`
}

// UpdateAccountRequest represents the request for updating an account
type UpdateAccountRequest struct {
	Name      *string      `json:"name"`
	Type      *AccountType `json:"type"`
	IsDefault *bool        `json:"is_default"`
}

// AccountSummary represents an account with its current balance
type AccountSummary struct {
	Account        *Account `json:"account"`
	CurrentBalance float64  `json:"current_balance"`
}
