package services

import (
	"Log45/budget/backend/db"
	"Log45/budget/backend/models"
	"context"
)

type UserService struct {
	users *db.UserRepository
}

func NewUserService(users *db.UserRepository) *UserService {
	return &UserService{
		users: users,
	}
}

func (s *UserService) GetByID(
	ctx context.Context,
	userID int64,
) (*models.User, error) {

	record, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:       record.ID,
		Username: record.Username,
		Email:    record.Email,
	}, nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error)
func (s *UserService) Update(models.User) error
func (s *UserService) Delete(models.User) error
