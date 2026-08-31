package db

import (
	"Log45/budget/backend/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrBudgetNotFound = errors.New("budget not found")

type BudgetRepository struct {
	pool *pgxpool.Pool
}

func NewBudgetRepository(pool *pgxpool.Pool) *BudgetRepository { return &BudgetRepository{pool: pool} }
func (r *BudgetRepository) Create(ctx context.Context, b *models.Budget) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, `INSERT INTO budgets (user_id,name,type,net_income,balance,start_date,end_date) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id,created_at,updated_at`, b.UserID, b.Name, b.Type, b.NetIncome, b.Balance, b.StartDate, b.EndDate).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return err
	}
	if err = r.replaceCategories(ctx, tx, b); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *BudgetRepository) List(ctx context.Context, userID int64) ([]models.Budget, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,user_id,name,type,net_income,balance,start_date,end_date,created_at,updated_at FROM budgets WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.Budget{}
	for rows.Next() {
		var b models.Budget
		if err := scanBudget(rows, &b); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}
func (r *BudgetRepository) Get(ctx context.Context, userID, id int64) (*models.Budget, error) {
	var b models.Budget
	err := scanBudget(r.pool.QueryRow(ctx, `SELECT id,user_id,name,type,net_income,balance,start_date,end_date,created_at,updated_at FROM budgets WHERE id=$1 AND user_id=$2`, id, userID), &b)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrBudgetNotFound
	}
	if err != nil {
		return nil, err
	}
	rows, err := r.pool.Query(ctx, `SELECT category_id,planned_amount FROM budget_categories WHERE budget_id=$1 ORDER BY category_id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c models.BudgetCategory
		if err := rows.Scan(&c.CategoryID, &c.PlannedAmount); err != nil {
			return nil, err
		}
		b.Categories = append(b.Categories, c)
	}
	return &b, rows.Err()
}
func (r *BudgetRepository) Update(ctx context.Context, b *models.Budget) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, `UPDATE budgets SET name=$1,type=$2,net_income=$3,balance=$4,start_date=$5,end_date=$6,updated_at=NOW() WHERE id=$7 AND user_id=$8 RETURNING created_at,updated_at`, b.Name, b.Type, b.NetIncome, b.Balance, b.StartDate, b.EndDate, b.ID, b.UserID).Scan(&b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrBudgetNotFound
	}
	if err != nil {
		return err
	}
	// A missing categories field represents a balance-only update (for example,
	// when an expense is linked to this budget). An explicit empty array clears
	// all allocations.
	if b.Categories != nil {
		if err = r.uncategorizeRemovedExpenses(ctx, tx, b); err != nil {
			return err
		}
		if err = r.replaceCategories(ctx, tx, b); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *BudgetRepository) uncategorizeRemovedExpenses(ctx context.Context, tx pgx.Tx, b *models.Budget) error {
	rows, err := tx.Query(ctx, `SELECT category_id FROM budget_categories WHERE budget_id=$1`, b.ID)
	if err != nil {
		return err
	}
	existing := []int64{}
	for rows.Next() {
		var categoryID int64
		if err := rows.Scan(&categoryID); err != nil {
			rows.Close()
			return err
		}
		existing = append(existing, categoryID)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	kept := make(map[int64]bool, len(b.Categories))
	for _, category := range b.Categories {
		kept[category.CategoryID] = true
	}
	for _, categoryID := range existing {
		if kept[categoryID] {
			continue
		}
		if _, err := tx.Exec(ctx, `UPDATE transactions SET category_id=NULL,updated_at=NOW() WHERE user_id=$1 AND budget_id=$2 AND category_id=$3`, b.UserID, b.ID, categoryID); err != nil {
			return err
		}
	}
	return nil
}

func (r *BudgetRepository) Delete(ctx context.Context, userID, id int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM budgets WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrBudgetNotFound
	}
	return nil
}
func (r *BudgetRepository) replaceCategories(ctx context.Context, tx pgx.Tx, b *models.Budget) error {
	if _, err := tx.Exec(ctx, `DELETE FROM budget_categories WHERE budget_id=$1`, b.ID); err != nil {
		return err
	}
	for _, c := range b.Categories {
		if _, err := tx.Exec(ctx, `INSERT INTO budget_categories (budget_id,category_id,planned_amount) VALUES ($1,$2,$3)`, b.ID, c.CategoryID, c.PlannedAmount); err != nil {
			return err
		}
	}
	return nil
}

type scanner interface{ Scan(...any) error }

func scanBudget(row scanner, b *models.Budget) error {
	return row.Scan(&b.ID, &b.UserID, &b.Name, &b.Type, &b.NetIncome, &b.Balance, &b.StartDate, &b.EndDate, &b.CreatedAt, &b.UpdatedAt)
}
