package services

import (
	"Log45/budget/backend/db"
	"Log45/budget/backend/models"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
)

var ErrInvalidTransaction = errors.New("amount must be positive, type valid, and date required")

type TransactionService struct {
	transactions *db.TransactionRepository
	categories   *db.CategoryRepository
}

func NewTransactionService(t *db.TransactionRepository, c *db.CategoryRepository) *TransactionService {
	return &TransactionService{transactions: t, categories: c}
}
func (s *TransactionService) List(ctx context.Context, userID int64) ([]models.Transaction, error) {
	return s.transactions.List(ctx, userID)
}
func (s *TransactionService) Create(ctx context.Context, userID int64, t models.Transaction) (*models.Transaction, error) {
	if err := s.validate(ctx, userID, &t); err != nil {
		return nil, err
	}
	t.UserID = userID
	if t.ID == "" {
		id, err := newID()
		if err != nil {
			return nil, err
		}
		t.ID = id
	}
	if err := s.transactions.Create(ctx, &t); err != nil {
		return nil, err
	}
	return &t, nil
}
func (s *TransactionService) Update(ctx context.Context, userID int64, t models.Transaction) (*models.Transaction, error) {
	if err := s.validate(ctx, userID, &t); err != nil {
		return nil, err
	}
	t.UserID = userID
	if err := s.transactions.Update(ctx, &t); err != nil {
		return nil, err
	}
	return &t, nil
}
func (s *TransactionService) Delete(ctx context.Context, userID int64, id string) error {
	return s.transactions.Delete(ctx, userID, id)
}
func (s *TransactionService) validate(ctx context.Context, userID int64, t *models.Transaction) error {
	if t.Amount <= 0 || t.Type < models.Income || t.Type > models.Transfer || t.Date.IsZero() {
		return ErrInvalidTransaction
	}
	if strings.TrimSpace(t.Description) == "" {
		return errors.New("description is required")
	}
	if t.CategoryID != nil {
		ok, err := s.categories.Accessible(ctx, userID, *t.CategoryID)
		if err != nil {
			return err
		}
		if !ok {
			return db.ErrCategoryNotFound
		}
	}
	return nil
}
func newID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
