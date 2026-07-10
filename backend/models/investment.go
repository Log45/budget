package models

type InvestmentType int

const (
	Stock InvestmentType = iota
	Bond
	IndexFund
	ETF
	Savings
)

type Investment struct {
	Name         string         `json:"name"`
	Amount       Money          `json:"amount"`
	Type         InvestmentType `json:"type"`
	ExpectedRate float64        `json:"expected_rate"` // annual interest rate as a decimal (e.g., 0.05 for 5%)
	ExpenseRatio float64        `json:"expense_ratio"` // annual expense ratio as a percentage (e.g., 0.01 for 0.01%)
	// TODO: add more fields as needed, such as investment type (stocks, bonds, etc.), risk level, etc.
}
