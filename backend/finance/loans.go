package finance

import (
	"Log45/budget/backend/models"
	"math"
)

func CalculateMonthlyPayment(loan models.Loan) models.Money {
	monthlyRate := loan.Rate / 12
	term := float64(loan.Term)
	if monthlyRate == 0 {
		return loan.Principal / models.Money(term)
	}
	// Using the formula for an amortizing loan payment
	payment := float64(loan.Principal) * (monthlyRate * math.Pow(1+monthlyRate, term)) / (math.Pow(1+monthlyRate, term) - 1)
	return models.Money(payment)
}

func GenerateSchedules(loan models.Loan) []models.Money {
	schedules := make([]models.Money, loan.Term)
	monthlyRate := loan.Rate / 12
	for i := 0; i < loan.Term; i++ {
		schedules[i] = MonthlyInterest(loan.Principal, monthlyRate)
	}
	return schedules
}
