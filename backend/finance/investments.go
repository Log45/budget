package finance

// might be weird, but HYSA will be considered an investment for now, if it gets complicated, split into savings file.
type InvestmentType int

const (
	Stock InvestmentType = iota
	Bond
	IndexFund
	ETF
	Savings
)

type Investment struct {
	Name         string
	Amount       Money
	ExpectedRate float64 // annual interest rate as a decimal (e.g., 0.05 for 5%)
	ExpenseRatio float64 // annual expense ratio as a percentage (e.g., 0.01 for 0.01%)
}
