package services

import (
	"context"
	"errors"
	"time"

	"Log45/budget/backend/db"
	"Log45/budget/backend/finance"
	"Log45/budget/backend/models"
)

var ErrInvalidLoan = errors.New("loan name, positive balance, non-negative rate, positive term, and start date are required")

type LoanService struct{ loans *db.LoanRepository }

func NewLoanService(loans *db.LoanRepository) *LoanService { return &LoanService{loans: loans} }

func (s *LoanService) Create(ctx context.Context, loan models.Loan) (*models.Loan, error) {
	if err := validateLoan(loan); err != nil {
		return nil, err
	}
	loan.CreatedAt, loan.UpdatedAt = time.Now(), time.Now()
	if err := s.loans.Create(ctx, &loan); err != nil {
		return nil, err
	}
	return &loan, nil
}
func (s *LoanService) List(ctx context.Context, userID int64) ([]models.Loan, error) {
	return s.loans.List(ctx, userID)
}
func (s *LoanService) Get(ctx context.Context, userID, id int64) (*models.Loan, error) {
	return s.loans.Get(ctx, userID, id)
}
func (s *LoanService) RecordPayment(ctx context.Context, userID, id int64, amount models.Money, paidAt time.Time) (*models.Loan, error) {
	if amount <= 0 {
		return nil, errors.New("payment amount must be positive")
	}
	if paidAt.IsZero() {
		paidAt = time.Now()
	}
	return s.loans.RecordPayment(ctx, userID, id, amount, paidAt)
}
func (s *LoanService) Schedule(ctx context.Context, userID, id int64) ([]models.LoanPayment, error) {
	loan, err := s.Get(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	return finance.GenerateSchedule(*loan, CalculateMonthlyPayment(*loan)), nil
}
func (s *LoanService) Analytics(ctx context.Context, userID, id int64) (models.LoanAnalytics, error) {
	loan, err := s.Get(ctx, userID, id)
	if err != nil {
		return models.LoanAnalytics{}, err
	}
	schedule := finance.GenerateSchedule(*loan, CalculateMonthlyPayment(*loan))
	analytics := models.LoanAnalytics{MonthlyPayment: CalculateMonthlyPayment(*loan), RemainingBalance: loan.CurrentBalance}
	for _, p := range schedule {
		analytics.RemainingInterest += p.Interest
		analytics.PayoffDate = p.DueDate
	}
	return analytics, nil
}
func validateLoan(loan models.Loan) error {
	if loan.Name == "" || loan.Principal <= 0 || loan.CurrentBalance <= 0 || loan.CurrentBalance > loan.Principal || loan.Rate < 0 || loan.Term <= 0 || loan.StartDate.IsZero() {
		return ErrInvalidLoan
	}
	return nil
}

func CalculateMonthlyPayment(loan models.Loan) models.Money {
	return finance.CalculateMonthlyPayment(loan)
}

func GenerateSchedules(loan models.Loan) []models.Money {
	return finance.GenerateSchedules(loan, CalculateMonthlyPayment(loan))
}
