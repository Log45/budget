package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrUserNotFound is returned when a user record cannot be found or is soft-deleted.
var ErrUserNotFound = errors.New("user not found")

// UserRecord represents a persisted user row including auth credentials.
type UserRecord struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// UserRepository provides PostgreSQL persistence for user records.
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository returns a UserRepository backed by the given connection pool.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Create inserts a new user and assigns the generated ID back to the record.
func (r *UserRepository) Create(ctx context.Context, user *UserRecord) error {
	const query = `
		INSERT INTO users (username, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetByID returns the active (non-deleted) user with the given ID.
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*UserRecord, error) {
	const query = `
		SELECT id, username, email, password_hash, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	return r.scanUser(ctx, query, id)
}

// GetByUsername returns the active user with the given username.
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*UserRecord, error) {
	const query = `
		SELECT id, username, email, password_hash, created_at, updated_at, deleted_at
		FROM users
		WHERE username = $1 AND deleted_at IS NULL
	`

	return r.scanUser(ctx, query, username)
}

// GetByEmail returns the active user with the given email address.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*UserRecord, error) {
	const query = `
		SELECT id, username, email, password_hash, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	return r.scanUser(ctx, query, email)
}

// Update modifies an existing user's username and email, refreshing updated_at.
func (r *UserRepository) Update(ctx context.Context, user *UserRecord) error {
	user.UpdatedAt = time.Now()

	const query = `
		UPDATE users
		SET username = $1, email = $2, updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	tag, err := r.pool.Exec(ctx, query, user.Username, user.Email, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Delete soft-deletes the user by setting deleted_at to the current time.
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE users
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

// scanUser executes a single-row query and maps the result into a UserRecord.
func (r *UserRepository) scanUser(ctx context.Context, query string, arg any) (*UserRecord, error) {
	var user UserRecord

	err := r.pool.QueryRow(ctx, query, arg).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
