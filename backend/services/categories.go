package services

import (
	"Log45/budget/backend/db"
	"Log45/budget/backend/models"
	"context"
	"errors"
	"strings"
)

var ErrInvalidCategory = errors.New("category name is required")

type CategoryService struct{ categories *db.CategoryRepository }

func NewCategoryService(categories *db.CategoryRepository) *CategoryService {
	return &CategoryService{categories: categories}
}
func (s *CategoryService) List(ctx context.Context, userID int64) ([]models.Category, error) {
	return s.categories.List(ctx, userID)
}
func (s *CategoryService) Create(ctx context.Context, userID int64, category models.Category) (*models.Category, error) {
	category.Name = strings.TrimSpace(category.Name)
	if category.Name == "" || category.TransactionType < models.Income || category.TransactionType > models.Transfer {
		return nil, ErrInvalidCategory
	}
	category.UserID = &userID
	if err := s.categories.Create(ctx, &category); err != nil {
		return nil, err
	}
	return &category, nil
}
