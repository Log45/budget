package models

import "time"

type Loan struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	Name           string    `json:"name"` // "Car Loan", "Mortgage"
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Principal      Money     `json:"principal"`
	CurrentBalance Money     `json:"current_balance"`
	Rate           float64   `json:"rate"`
	Term           int       `json:"term"` // months
	StartDate      time.Time `json:"start_date"`
}
