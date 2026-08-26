package models

import "time"

type BudgetType int

const (
	YearlyBudget BudgetType = iota
	MonthlyBudget
	WeeklyBudget
	DailyBudget
)

type Budget struct {
	ID         int64            `json:"id"`
	UserID     int64            `json:"user_id"`
	Name       string           `json:"name"`
	Type       BudgetType       `json:"type"`
	NetIncome  Money            `json:"net_income"`
	Balance    Money            `json:"balance"`
	StartDate  *time.Time       `json:"start_date,omitempty"`
	EndDate    *time.Time       `json:"end_date,omitempty"`
	Categories []BudgetCategory `json:"categories,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

// BudgetCategory plans an amount for a category within one budget scenario.
// It does not own transactions; transactions are linked to both a budget and a category.
type BudgetCategory struct {
	CategoryID    int64 `json:"category_id"`
	PlannedAmount Money `json:"planned_amount"`
}
