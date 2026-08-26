package finance

import (
	"Log45/budget/backend/models"
	"math"
)

func CalculateMonthlyPayment(loan models.Loan) models.Money {
	if loan.Term <= 0 || loan.CurrentBalance <= 0 {
		return 0
	}
	monthlyRate := loan.Rate / 12
	term := float64(loan.Term)
	if monthlyRate == 0 {
		return loan.CurrentBalance / models.Money(term)
	}
	// Using the formula for an amortizing loan payment
	payment := float64(loan.CurrentBalance) * (monthlyRate * math.Pow(1+monthlyRate, term)) / (math.Pow(1+monthlyRate, term) - 1)
	return models.Money(math.Ceil(payment))
}

// GenerateSchedule creates an amortization schedule starting from the current balance.
func GenerateSchedule(loan models.Loan, monthlyPayment models.Money) []models.LoanPayment {
	if loan.Term <= 0 || loan.CurrentBalance <= 0 || monthlyPayment <= 0 {
		return []models.LoanPayment{}
	}
	balance := loan.CurrentBalance
	monthlyRate := loan.Rate / 12
	schedule := make([]models.LoanPayment, 0, loan.Term)
	for i := 0; i < loan.Term && balance > 0; i++ {
		interest := models.Money(math.Round(float64(balance) * monthlyRate))
		principal := monthlyPayment - interest
		payment := monthlyPayment
		if principal <= 0 {
			break
		}
		if principal >= balance {
			principal = balance
			payment = principal + interest
		}
		balance -= principal
		schedule = append(schedule, models.LoanPayment{Number: i + 1, DueDate: loan.StartDate.AddDate(0, i+1, 0), Payment: payment, Principal: principal, Interest: interest, Balance: balance})
	}
	return schedule
}

// GenerateSchedules is retained for callers that only need each payment's interest.
func GenerateSchedules(loan models.Loan, monthlyPayment models.Money) []models.Money {
	schedule := GenerateSchedule(loan, monthlyPayment)
	result := make([]models.Money, len(schedule))
	for i, payment := range schedule {
		result[i] = payment.Interest
	}
	return result
}
