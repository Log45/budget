package finance

import (
	"Log45/budget/backend/models"
	"math"
)

func CalculateMonthlyPayment(loan models.Loan) models.Money {
	monthlyRate := loan.Rate / float64(loan.Term)
	term := float64(loan.Term)
	if monthlyRate == 0 {
		return loan.Principal / models.Money(term)
	}
	// Using the formula for an amortizing loan payment
	payment := float64(loan.Principal) * (monthlyRate * math.Pow(1+monthlyRate, term)) / (math.Pow(1+monthlyRate, term) - 1)
	return models.Money(payment)
}

// TODO: Change this to use the actual balance to generate an amortization schedule. Needs per-payment principal/interest/balance breakdown
func GenerateSchedules(loan models.Loan, monthlyPayment models.Money) []models.Money {
	schedules := make([]models.Money, loan.Term)
	monthlyRate := loan.Rate / 12
	for i := 0; i < loan.Term; i++ {
		schedules[i] = MonthlyInterest(loan.Principal, monthlyRate)
	}
	return schedules
}
