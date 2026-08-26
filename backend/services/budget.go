package services

import (
	"Log45/budget/backend/db"
	"Log45/budget/backend/models"
	"context"
	"errors"
	"strings"
)

var ErrInvalidBudget = errors.New("budget name, valid type, non-negative income, and valid category allocations are required")

type BudgetService struct {
	budgets    *db.BudgetRepository
	categories *db.CategoryRepository
}

func NewBudgetService(b *db.BudgetRepository, c *db.CategoryRepository) *BudgetService {
	return &BudgetService{budgets: b, categories: c}
}
func (s *BudgetService) Create(ctx context.Context, userID int64, b models.Budget) (*models.Budget, error) {
	b.UserID = userID
	if err := s.validate(ctx, userID, &b); err != nil {
		return nil, err
	}
	if err := s.budgets.Create(ctx, &b); err != nil {
		return nil, err
	}
	return &b, nil
}
func (s *BudgetService) List(ctx context.Context, userID int64) ([]models.Budget, error) {
	return s.budgets.List(ctx, userID)
}
func (s *BudgetService) Get(ctx context.Context, userID, id int64) (*models.Budget, error) {
	return s.budgets.Get(ctx, userID, id)
}
func (s *BudgetService) Update(ctx context.Context, userID int64, b models.Budget) (*models.Budget, error) {
	b.UserID = userID
	if err := s.validate(ctx, userID, &b); err != nil {
		return nil, err
	}
	if err := s.budgets.Update(ctx, &b); err != nil {
		return nil, err
	}
	return &b, nil
}
func (s *BudgetService) Delete(ctx context.Context, userID, id int64) error {
	return s.budgets.Delete(ctx, userID, id)
}
func (s *BudgetService) validate(ctx context.Context, userID int64, b *models.Budget) error {
	b.Name = strings.TrimSpace(b.Name)
	if b.Name == "" || b.Type < models.YearlyBudget || b.Type > models.DailyBudget || b.NetIncome < 0 || (b.StartDate != nil && b.EndDate != nil && b.EndDate.Before(*b.StartDate)) {
		return ErrInvalidBudget
	}
	seen := map[int64]bool{}
	for _, c := range b.Categories {
		if c.CategoryID <= 0 || c.PlannedAmount < 0 || seen[c.CategoryID] {
			return ErrInvalidBudget
		}
		seen[c.CategoryID] = true
		ok, err := s.categories.Accessible(ctx, userID, c.CategoryID)
		if err != nil {
			return err
		}
		if !ok {
			return db.ErrCategoryNotFound
		}
	}
	return nil
}
