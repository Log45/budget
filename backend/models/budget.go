package models

type BudgetType int

const (
	YearlyBudget BudgetType = iota
	MonthlyBudget
	WeeklyBudget
	DailyBudget
)

type Budget struct {
	ID        int64
	Type      BudgetType
	NetIncome Money
	Balance   Money
	Expenses  []Transaction
}
