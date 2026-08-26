package models

import "time"

type Property struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Name          string    `json:"name"`
	Address       string    `json:"address"`
	Type          string    `json:"type"`
	PurchasePrice *Money    `json:"purchase_price,omitempty"`
	CurrentValue  *Money    `json:"current_value,omitempty"`
	LoanID        *int64    `json:"loan_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PropertyAnalytics struct {
	Income   Money  `json:"income"`
	Expenses Money  `json:"expenses"`
	CashFlow Money  `json:"cash_flow"`
	Equity   *Money `json:"equity,omitempty"`
}
