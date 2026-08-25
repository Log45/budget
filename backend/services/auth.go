package services

import (
	"Log45/budget/backend/db"

	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ErrUsernameExists is returned when registration is attempted with a taken username.
var ErrUsernameExists = errors.New("username already exists")

// ErrEmailExists is returned when registration is attempted with a taken email.
var ErrEmailExists = errors.New("email already exists")

// AuthService handles user registration, login, and JWT token generation.
type AuthService struct {
	users     *db.UserRepository
	jwtSecret []byte
}

// Claims holds the JWT payload for authenticated requests.
type Claims struct {
	UserID int64 `json:"user_id"`

	jwt.RegisteredClaims
}

// NewAuthService returns an AuthService using the given user repository and JWT secret.
func NewAuthService(
	users *db.UserRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		users:     users,
		jwtSecret: []byte(jwtSecret),
	}
}

// Register creates a new user account after checking username and email uniqueness.
func (s *AuthService) Register(
	ctx context.Context,
	username string,
	email string,
	password string,
) error {
	// check if the username exists
	if _, err := s.users.GetByUsername(ctx, username); err == nil {
		return ErrUsernameExists
	}

	// check if the email exists
	if _, err := s.users.GetByEmail(ctx, email); err == nil {
		return ErrEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := &db.UserRecord{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := s.users.Create(ctx, user); err != nil {
		return err
	}
	return nil
}

// GenerateToken creates a signed JWT for the given user ID, valid for 24 hours.
func (s *AuthService) GenerateToken(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "budget-app",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// Login verifies credentials and returns a JWT for the authenticated user.
func (s *AuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (string, error) {
	user, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", err
	}
	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}
