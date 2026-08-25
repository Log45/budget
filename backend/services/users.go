package services

import (
	"context"

	"Log45/budget/backend/db"
	"Log45/budget/backend/models"
)

// UserService provides business logic for user records.
type UserService struct {
	users *db.UserRepository
}

// NewUserService returns a UserService backed by the given repository.
func NewUserService(users *db.UserRepository) *UserService {
	return &UserService{users: users}
}

// GetByID returns the public user profile for the given ID.
func (s *UserService) GetByID(ctx context.Context, userID int64) (*models.User, error) {
	record, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return recordToUser(record), nil
}

// GetByUsername returns the public user profile for the given username.
func (s *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	record, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return recordToUser(record), nil
}

// Update modifies an existing user's username and email.
func (s *UserService) Update(ctx context.Context, user models.User) error {
	return s.users.Update(ctx, &db.UserRecord{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}

// Delete soft-deletes the user with the given ID.
func (s *UserService) Delete(ctx context.Context, userID int64) error {
	return s.users.Delete(ctx, userID)
}

// recordToUser maps a database record to the public API user model.
func recordToUser(record *db.UserRecord) *models.User {
	return &models.User{
		ID:       record.ID,
		Username: record.Username,
		Email:    record.Email,
	}
}
