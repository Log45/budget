package services

import (
	"Log45/budget/backend/finance"
	"Log45/budget/backend/models"
)

func CalculateMonthlyPayment(loan models.Loan) models.Money {
	return finance.CalculateMonthlyPayment(loan)
}

func GenerateSchedules(loan models.Loan) []models.Money {
	return finance.GenerateSchedules(loan)
}
