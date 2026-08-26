package db

import (
	"Log45/budget/backend/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var ErrLoanNotFound = errors.New("loan not found")

type LoanRepository struct {
	pool *pgxpool.Pool
}

func NewLoanRepository(pool *pgxpool.Pool) *LoanRepository { return &LoanRepository{pool: pool} }
func (r *LoanRepository) Create(ctx context.Context, loan *models.Loan) error {
	return r.pool.QueryRow(ctx, `INSERT INTO loans (user_id,name,principal,current_balance,rate,term,start_date,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`, loan.UserID, loan.Name, loan.Principal, loan.CurrentBalance, loan.Rate, loan.Term, loan.StartDate, loan.CreatedAt, loan.UpdatedAt).Scan(&loan.ID)
}
func (r *LoanRepository) List(ctx context.Context, userID int64) ([]models.Loan, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,user_id,name,principal,current_balance,rate,term,start_date,created_at,updated_at FROM loans WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	loans := []models.Loan{}
	for rows.Next() {
		var l models.Loan
		if err := rows.Scan(&l.ID, &l.UserID, &l.Name, &l.Principal, &l.CurrentBalance, &l.Rate, &l.Term, &l.StartDate, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		loans = append(loans, l)
	}
	return loans, rows.Err()
}
func (r *LoanRepository) Get(ctx context.Context, userID, id int64) (*models.Loan, error) {
	var l models.Loan
	err := r.pool.QueryRow(ctx, `SELECT id,user_id,name,principal,current_balance,rate,term,start_date,created_at,updated_at FROM loans WHERE id=$1 AND user_id=$2`, id, userID).Scan(&l.ID, &l.UserID, &l.Name, &l.Principal, &l.CurrentBalance, &l.Rate, &l.Term, &l.StartDate, &l.CreatedAt, &l.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrLoanNotFound
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}
func (r *LoanRepository) RecordPayment(ctx context.Context, userID, id int64, amount models.Money, paidAt time.Time) (*models.Loan, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var loan models.Loan
	err = tx.QueryRow(ctx, `SELECT id,user_id,name,principal,current_balance,rate,term,start_date,created_at,updated_at FROM loans WHERE id=$1 AND user_id=$2 FOR UPDATE`, id, userID).Scan(&loan.ID, &loan.UserID, &loan.Name, &loan.Principal, &loan.CurrentBalance, &loan.Rate, &loan.Term, &loan.StartDate, &loan.CreatedAt, &loan.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrLoanNotFound
	}
	if err != nil {
		return nil, err
	}
	applied := amount
	if applied > loan.CurrentBalance {
		applied = loan.CurrentBalance
	}
	loan.CurrentBalance -= applied
	err = tx.QueryRow(ctx, `UPDATE loans SET current_balance=$1,updated_at=NOW() WHERE id=$2 RETURNING updated_at`, loan.CurrentBalance, id).Scan(&loan.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(ctx, `INSERT INTO loan_payments (loan_id,amount,paid_at) VALUES ($1,$2,$3)`, id, applied, paidAt)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &loan, nil
}
