package models

import "time"

type TransactionType int

const (
	Income TransactionType = iota
	Expense
	Transfer
)

type Transaction struct {
	ID          string          `json:"id"`
	UserID      int64           `json:"user_id"`
	BudgetID    *int64          `json:"budget_id,omitempty"`
	CategoryID  *int64          `json:"category_id,omitempty"`
	PropertyID  *int64          `json:"property_id,omitempty"`
	AccountID   *int64          `json:"account_id,omitempty"`
	Amount      Money           `json:"amount"`
	Description string          `json:"description"`
	Type        TransactionType `json:"type"`
	Source      string          `json:"source"`      // e.g., "Bank", "Credit Card", "Cash"
	Destination string          `json:"destination"` // e.g., "Bank", "Credit Card", "Cash" TODO: make this a pointer to an account struct, or maybe just an account ID string
	Date        time.Time       `json:"date"`        // ISO 8601 format (YYYY-MM-DD)
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
