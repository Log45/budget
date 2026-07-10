package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRecord struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *UserRecord,
) error

func (r *UserRepository) GetByID(
	ctx context.Context,
	id int64,
) (*UserRecord, error)

func (r *UserRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*UserRecord, error)

func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*UserRecord, error)

func (r *UserRepository) Update(
	ctx context.Context,
	user *UserRecord,
) error

func (r *UserRepository) Delete(
	ctx context.Context,
	id int64,
) error
