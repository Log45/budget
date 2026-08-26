package db

import (
	"Log45/budget/backend/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTransactionNotFound = errors.New("transaction not found")

type TransactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{pool: pool}
}
func (r *TransactionRepository) Create(ctx context.Context, t *models.Transaction) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := tx.QueryRow(ctx, `INSERT INTO transactions (id,user_id,budget_id,category_id,property_id,account_id,amount,description,type,source,destination,date) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING created_at,updated_at`, t.ID, t.UserID, t.BudgetID, t.CategoryID, t.PropertyID, t.AccountID, t.Amount, t.Description, t.Type, t.Source, t.Destination, t.Date).Scan(&t.CreatedAt, &t.UpdatedAt); err != nil {
		return err
	}
	if t.Type == models.Expense && t.BudgetID != nil && t.CategoryID != nil {
		tag, err := tx.Exec(ctx, `UPDATE budget_categories SET planned_amount = planned_amount + $1 WHERE budget_id = $2 AND category_id = $3`, t.Amount, *t.BudgetID, *t.CategoryID)
		if err != nil {
			return err
		}
		if tag.RowsAffected() > 0 {
			if _, err := tx.Exec(ctx, `UPDATE budgets SET balance = balance - $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`, t.Amount, *t.BudgetID, t.UserID); err != nil {
				return err
			}
		}
	}
	return tx.Commit(ctx)
}
func (r *TransactionRepository) List(ctx context.Context, userID int64) ([]models.Transaction, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,user_id,budget_id,category_id,property_id,account_id,amount,description,type,source,destination,date,created_at,updated_at FROM transactions WHERE user_id=$1 ORDER BY date DESC,created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []models.Transaction{}
	for rows.Next() {
		t, err := scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}
func (r *TransactionRepository) Get(ctx context.Context, userID int64, id string) (*models.Transaction, error) {
	t, err := scanTransaction(r.pool.QueryRow(ctx, `SELECT id,user_id,budget_id,category_id,property_id,account_id,amount,description,type,source,destination,date,created_at,updated_at FROM transactions WHERE id=$1 AND user_id=$2`, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTransactionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}
func (r *TransactionRepository) Update(ctx context.Context, t *models.Transaction) error {
	err := r.pool.QueryRow(ctx, `UPDATE transactions SET budget_id=$1,category_id=$2,property_id=$3,account_id=$4,amount=$5,description=$6,type=$7,source=$8,destination=$9,date=$10,updated_at=NOW() WHERE id=$11 AND user_id=$12 RETURNING created_at,updated_at`, t.BudgetID, t.CategoryID, t.PropertyID, t.AccountID, t.Amount, t.Description, t.Type, t.Source, t.Destination, t.Date, t.ID, t.UserID).Scan(&t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrTransactionNotFound
	}
	return err
}
func (r *TransactionRepository) Delete(ctx context.Context, userID int64, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM transactions WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrTransactionNotFound
	}
	return nil
}

type rowScanner interface{ Scan(...any) error }

func scanTransaction(row rowScanner) (models.Transaction, error) {
	var t models.Transaction
	err := row.Scan(&t.ID, &t.UserID, &t.BudgetID, &t.CategoryID, &t.PropertyID, &t.AccountID, &t.Amount, &t.Description, &t.Type, &t.Source, &t.Destination, &t.Date, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}
