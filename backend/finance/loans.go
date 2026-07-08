package finance

import (
	"math"
)

type Loan struct {
	Principal Money
	Rate      float64 // Annual interest rate as a decimal (e.g., 0.05 for 5%)
	Term      int     // Term in months
}

type Payment struct {
	Amount Money
	Date   string // YYYY-MM-DD format
}

func CalculateMonthlyPayment(loan Loan) Money {
	monthlyRate := loan.Rate / 12
	term := float64(loan.Term)
	if monthlyRate == 0 {
		return loan.Principal / Money(term)
	}
	// Using the formula for an amortizing loan payment
	payment := float64(loan.Principal) * (monthlyRate * math.Pow(1+monthlyRate, term)) / (math.Pow(1+monthlyRate, term) - 1)
	return Money(payment)
}

func GenerateSchedules(loan Loan) []Money {
	schedules := make([]Money, loan.Term)
	monthlyRate := loan.Rate / 12
	for i := 0; i < loan.Term; i++ {
		schedules[i] = MonthlyInterest(loan.Principal, monthlyRate)
	}
	return schedules
}
