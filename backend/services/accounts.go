package services

import (
	"Log45/budget/backend/db"
	"Log45/budget/backend/models"
	"context"
	"errors"
	"strings"
)

var ErrInvalidAccount = errors.New("account name and valid account type are required; credit limit cannot be negative")

type AccountService struct{ accounts *db.AccountRepository }

func NewAccountService(accounts *db.AccountRepository) *AccountService {
	return &AccountService{accounts: accounts}
}
func (s *AccountService) Create(ctx context.Context, userID int64, a models.Account) (*models.Account, error) {
	a.UserID = userID
	if err := validateAccount(&a); err != nil {
		return nil, err
	}
	if err := s.accounts.Create(ctx, &a); err != nil {
		return nil, err
	}
	return &a, nil
}
func (s *AccountService) List(ctx context.Context, userID int64) ([]models.Account, error) {
	return s.accounts.List(ctx, userID)
}
func (s *AccountService) Get(ctx context.Context, userID, id int64) (*models.Account, error) {
	return s.accounts.Get(ctx, userID, id)
}
func (s *AccountService) Update(ctx context.Context, userID int64, a models.Account) (*models.Account, error) {
	a.UserID = userID
	if err := validateAccount(&a); err != nil {
		return nil, err
	}
	if err := s.accounts.Update(ctx, &a); err != nil {
		return nil, err
	}
	return &a, nil
}
func (s *AccountService) Delete(ctx context.Context, userID, id int64) error {
	return s.accounts.Delete(ctx, userID, id)
}

func validateAccount(a *models.Account) error {
	a.Name = strings.TrimSpace(a.Name)
	valid := map[models.AccountType]bool{models.SavingsAccount: true, models.CheckingAccount: true, models.HybridAccount: true, models.IndividualBrokerage: true, models.RothIRAAccount: true, models.IRAAccount: true, models.Account401K: true, models.Account403B: true, models.Account457B: true, models.CreditCardAccount: true}
	if a.Name == "" || !valid[a.AccountType] || (a.Limit != nil && *a.Limit < 0) {
		return ErrInvalidAccount
	}
	if a.AccountType != models.CreditCardAccount {
		a.Limit = nil
	}
	return nil
}
