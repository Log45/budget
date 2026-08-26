package models

import "time"

type Loan struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	Name           string    `json:"name"` // "Car Loan", "Mortgage"
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Principal      Money     `json:"principal"`
	CurrentBalance Money     `json:"current_balance"`
	Rate           float64   `json:"rate"` // annual rate as a decimal, e.g. 0.065 for 6.5%
	Term           int       `json:"term"` // months
	StartDate      time.Time `json:"start_date"`
}

// LoanPayment is one projected or recorded payment. All money values are cents.
type LoanPayment struct {
	Number    int       `json:"number"`
	DueDate   time.Time `json:"due_date"`
	Payment   Money     `json:"payment"`
	Principal Money     `json:"principal"`
	Interest  Money     `json:"interest"`
	Balance   Money     `json:"balance"`
}

type LoanAnalytics struct {
	MonthlyPayment    Money     `json:"monthly_payment"`
	RemainingBalance  Money     `json:"remaining_balance"`
	RemainingInterest Money     `json:"remaining_interest"`
	PayoffDate        time.Time `json:"payoff_date"`
}
