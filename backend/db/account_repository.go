package db

import (
	"Log45/budget/backend/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAccountNotFound = errors.New("account not found")

type AccountRepository struct{ pool *pgxpool.Pool }

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{pool: pool}
}
func (r *AccountRepository) Create(ctx context.Context, a *models.Account) error {
	return r.pool.QueryRow(ctx, `INSERT INTO accounts (user_id,name,account_type,balance,credit_limit) VALUES ($1,$2,$3,$4,$5) RETURNING id,created_at,updated_at`, a.UserID, a.Name, a.AccountType, a.Balance, a.Limit).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}
func (r *AccountRepository) List(ctx context.Context, userID int64) ([]models.Account, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,user_id,name,account_type,balance,credit_limit,created_at,updated_at FROM accounts WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.Account{}
	for rows.Next() {
		var a models.Account
		if err := scanAccount(rows, &a); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, rows.Err()
}
func (r *AccountRepository) Get(ctx context.Context, userID, id int64) (*models.Account, error) {
	var a models.Account
	err := scanAccount(r.pool.QueryRow(ctx, `SELECT id,user_id,name,account_type,balance,credit_limit,created_at,updated_at FROM accounts WHERE id=$1 AND user_id=$2`, id, userID), &a)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrAccountNotFound
	}
	return &a, err
}
func (r *AccountRepository) Update(ctx context.Context, a *models.Account) error {
	err := r.pool.QueryRow(ctx, `UPDATE accounts SET name=$1,account_type=$2,balance=$3,credit_limit=$4,updated_at=NOW() WHERE id=$5 AND user_id=$6 RETURNING created_at,updated_at`, a.Name, a.AccountType, a.Balance, a.Limit, a.ID, a.UserID).Scan(&a.CreatedAt, &a.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrAccountNotFound
	}
	return err
}
func (r *AccountRepository) Delete(ctx context.Context, userID, id int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM accounts WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrAccountNotFound
	}
	return nil
}
func scanAccount(row scanner, a *models.Account) error {
	return row.Scan(&a.ID, &a.UserID, &a.Name, &a.AccountType, &a.Balance, &a.Limit, &a.CreatedAt, &a.UpdatedAt)
}
