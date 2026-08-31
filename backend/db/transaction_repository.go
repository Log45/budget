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
	if err := tx.QueryRow(ctx, `INSERT INTO transactions (id,user_id,budget_id,category_id,property_id,account_id,amount,description,type,source,destination,date,recurring,recurring_type) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING created_at,updated_at`, t.ID, t.UserID, t.BudgetID, t.CategoryID, t.PropertyID, t.AccountID, t.Amount, t.Description, t.Type, t.Source, t.Destination, t.Date, t.Recurring, t.RecurringType).Scan(&t.CreatedAt, &t.UpdatedAt); err != nil {
		return err
	}
	if err := adjustBudget(ctx, tx, t, 1); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *TransactionRepository) List(ctx context.Context, userID int64) ([]models.Transaction, error) {
	rows, err := r.pool.Query(ctx, transactionSelect+` WHERE t.user_id=$1 ORDER BY t.date DESC,t.created_at DESC`, userID)
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
	t, err := scanTransaction(r.pool.QueryRow(ctx, transactionSelect+` WHERE t.id=$1 AND t.user_id=$2`, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTransactionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}
func (r *TransactionRepository) Update(ctx context.Context, t *models.Transaction) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	old, err := scanTransaction(tx.QueryRow(ctx, transactionSelect+` WHERE t.id=$1 AND t.user_id=$2 FOR UPDATE OF t`, t.ID, t.UserID))
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrTransactionNotFound
	}
	if err != nil {
		return err
	}
	if err = adjustBudget(ctx, tx, &old, -1); err != nil {
		return err
	}
	err = tx.QueryRow(ctx, `UPDATE transactions SET budget_id=$1,category_id=$2,property_id=$3,account_id=$4,amount=$5,description=$6,type=$7,source=$8,destination=$9,date=$10,recurring=$11,recurring_type=$12,updated_at=NOW() WHERE id=$13 AND user_id=$14 RETURNING created_at,updated_at`, t.BudgetID, t.CategoryID, t.PropertyID, t.AccountID, t.Amount, t.Description, t.Type, t.Source, t.Destination, t.Date, t.Recurring, t.RecurringType, t.ID, t.UserID).Scan(&t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return err
	}
	if err = adjustBudget(ctx, tx, t, 1); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *TransactionRepository) Delete(ctx context.Context, userID int64, id string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	old, err := scanTransaction(tx.QueryRow(ctx, transactionSelect+` WHERE t.id=$1 AND t.user_id=$2 FOR UPDATE OF t`, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrTransactionNotFound
	}
	if err != nil {
		return err
	}
	if err = adjustBudget(ctx, tx, &old, -1); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM transactions WHERE id=$1 AND user_id=$2`, id, userID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

const transactionSelect = `SELECT t.id,t.user_id,t.budget_id,t.category_id,t.property_id,t.account_id,t.amount,t.description,t.type,t.source,t.destination,t.date,t.created_at,t.updated_at,t.recurring,t.recurring_type,CASE WHEN b.net_income > 0 THEN (t.amount::double precision / b.net_income::double precision) * 100 ELSE 0 END FROM transactions t LEFT JOIN budgets b ON b.id=t.budget_id`

func adjustBudget(ctx context.Context, tx pgx.Tx, t *models.Transaction, direction int64) error {
	if t.Type != models.Expense || t.BudgetID == nil || t.CategoryID == nil {
		return nil
	}
	requestedDelta := int64(t.Amount) * direction
	var current models.Money
	inserted := false
	err := tx.QueryRow(ctx, `SELECT planned_amount FROM budget_categories WHERE budget_id=$1 AND category_id=$2 FOR UPDATE`, *t.BudgetID, *t.CategoryID).Scan(&current)
	if errors.Is(err, pgx.ErrNoRows) {
		if direction < 0 {
			// Older expenses may be linked to a budget without ever having been
			// included in its category allocation. There is nothing to reverse.
			return nil
		}
		if _, err = tx.Exec(ctx, `INSERT INTO budget_categories (budget_id,category_id,planned_amount) VALUES ($1,$2,$3)`, *t.BudgetID, *t.CategoryID, t.Amount); err != nil {
			return err
		}
		inserted = true
		current = 0
	} else if err != nil {
		return err
	}

	newAmount := int64(current) + requestedDelta
	if newAmount < 0 {
		newAmount = 0
	}
	appliedDelta := newAmount - int64(current)
	if inserted {
		// The no-row branch already inserted the requested amount.
		appliedDelta = requestedDelta
	} else {
		if _, err = tx.Exec(ctx, `UPDATE budget_categories SET planned_amount=$1 WHERE budget_id=$2 AND category_id=$3`, newAmount, *t.BudgetID, *t.CategoryID); err != nil {
			return err
		}
	}
	_, err = tx.Exec(ctx, `UPDATE budgets SET balance=balance-$1,updated_at=NOW() WHERE id=$2 AND user_id=$3`, appliedDelta, *t.BudgetID, t.UserID)
	return err
}

type rowScanner interface{ Scan(...any) error }

func scanTransaction(row rowScanner) (models.Transaction, error) {
	var t models.Transaction
	err := row.Scan(&t.ID, &t.UserID, &t.BudgetID, &t.CategoryID, &t.PropertyID, &t.AccountID, &t.Amount, &t.Description, &t.Type, &t.Source, &t.Destination, &t.Date, &t.CreatedAt, &t.UpdatedAt, &t.Recurring, &t.RecurringType, &t.BudgetRatio)
	return t, err
}
