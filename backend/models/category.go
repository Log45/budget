package models

import "time"

// Category is either a shared default (UserID nil) or owned by one user.
type Category struct {
	ID              int64           `json:"id"`
	UserID          *int64          `json:"user_id,omitempty"`
	Name            string          `json:"name"`
	TransactionType TransactionType `json:"transaction_type"`
	IsDefault       bool            `json:"is_default"`
	CreatedAt       time.Time       `json:"created_at"`
}
