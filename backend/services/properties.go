package services

import (
	"Log45/budget/backend/db"
	"Log45/budget/backend/models"
	"context"
	"errors"
	"strings"
)

var ErrInvalidProperty = errors.New("property name and type are required; values cannot be negative")

type PropertyService struct {
	properties *db.PropertyRepository
	loans      *db.LoanRepository
}

func NewPropertyService(p *db.PropertyRepository, l *db.LoanRepository) *PropertyService {
	return &PropertyService{properties: p, loans: l}
}
func (s *PropertyService) Create(ctx context.Context, userID int64, p models.Property) (*models.Property, error) {
	p.UserID = userID
	if err := s.validate(ctx, userID, &p); err != nil {
		return nil, err
	}
	if err := s.properties.Create(ctx, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
func (s *PropertyService) List(ctx context.Context, userID int64) ([]models.Property, error) {
	return s.properties.List(ctx, userID)
}
func (s *PropertyService) Get(ctx context.Context, userID, id int64) (*models.Property, error) {
	return s.properties.Get(ctx, userID, id)
}
func (s *PropertyService) Update(ctx context.Context, userID int64, p models.Property) (*models.Property, error) {
	p.UserID = userID
	if err := s.validate(ctx, userID, &p); err != nil {
		return nil, err
	}
	if err := s.properties.Update(ctx, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
func (s *PropertyService) Delete(ctx context.Context, userID, id int64) error {
	return s.properties.Delete(ctx, userID, id)
}
func (s *PropertyService) Analytics(ctx context.Context, userID, id int64) (models.PropertyAnalytics, error) {
	return s.properties.Analytics(ctx, userID, id)
}
func (s *PropertyService) validate(ctx context.Context, userID int64, p *models.Property) error {
	p.Name = strings.TrimSpace(p.Name)
	p.Type = strings.TrimSpace(p.Type)
	if p.Name == "" || p.Type == "" || (p.PurchasePrice != nil && *p.PurchasePrice < 0) || (p.CurrentValue != nil && *p.CurrentValue < 0) {
		return ErrInvalidProperty
	}
	if p.LoanID != nil {
		if _, err := s.loans.Get(ctx, userID, *p.LoanID); err != nil {
			return err
		}
	}
	return nil
}
