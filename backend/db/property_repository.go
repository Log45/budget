package db

import (
	"Log45/budget/backend/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrPropertyNotFound = errors.New("property not found")

type PropertyRepository struct{ pool *pgxpool.Pool }

func NewPropertyRepository(pool *pgxpool.Pool) *PropertyRepository {
	return &PropertyRepository{pool: pool}
}
func (r *PropertyRepository) Create(ctx context.Context, p *models.Property) error {
	return r.pool.QueryRow(ctx, `INSERT INTO properties (user_id,name,address,type,purchase_price,current_value,loan_id) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id,created_at,updated_at`, p.UserID, p.Name, p.Address, p.Type, p.PurchasePrice, p.CurrentValue, p.LoanID).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}
func (r *PropertyRepository) List(ctx context.Context, userID int64) ([]models.Property, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,user_id,name,address,type,purchase_price,current_value,loan_id,created_at,updated_at FROM properties WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.Property{}
	for rows.Next() {
		var p models.Property
		if err := scanProperty(rows, &p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
func (r *PropertyRepository) Get(ctx context.Context, userID, id int64) (*models.Property, error) {
	var p models.Property
	err := scanProperty(r.pool.QueryRow(ctx, `SELECT id,user_id,name,address,type,purchase_price,current_value,loan_id,created_at,updated_at FROM properties WHERE id=$1 AND user_id=$2`, id, userID), &p)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPropertyNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}
func (r *PropertyRepository) Update(ctx context.Context, p *models.Property) error {
	err := r.pool.QueryRow(ctx, `UPDATE properties SET name=$1,address=$2,type=$3,purchase_price=$4,current_value=$5,loan_id=$6,updated_at=NOW() WHERE id=$7 AND user_id=$8 RETURNING created_at,updated_at`, p.Name, p.Address, p.Type, p.PurchasePrice, p.CurrentValue, p.LoanID, p.ID, p.UserID).Scan(&p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrPropertyNotFound
	}
	return err
}
func (r *PropertyRepository) Delete(ctx context.Context, userID, id int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM properties WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrPropertyNotFound
	}
	return nil
}
func (r *PropertyRepository) Analytics(ctx context.Context, userID, id int64) (models.PropertyAnalytics, error) {
	p, err := r.Get(ctx, userID, id)
	if err != nil {
		return models.PropertyAnalytics{}, err
	}
	var a models.PropertyAnalytics
	err = r.pool.QueryRow(ctx, `SELECT COALESCE(SUM(amount) FILTER (WHERE type=0),0),COALESCE(SUM(amount) FILTER (WHERE type=1),0) FROM transactions WHERE user_id=$1 AND property_id=$2`, userID, id).Scan(&a.Income, &a.Expenses)
	if err != nil {
		return a, err
	}
	a.CashFlow = a.Income - a.Expenses
	if p.CurrentValue != nil {
		var balance *models.Money
		if p.LoanID != nil {
			err = r.pool.QueryRow(ctx, `SELECT current_balance FROM loans WHERE id=$1 AND user_id=$2`, *p.LoanID, userID).Scan(&balance)
			if errors.Is(err, pgx.ErrNoRows) {
				balance = nil
			} else if err != nil {
				return a, err
			}
		}
		value := *p.CurrentValue
		if balance != nil {
			value -= *balance
		}
		a.Equity = &value
	}
	return a, nil
}
func scanProperty(row scanner, p *models.Property) error {
	return row.Scan(&p.ID, &p.UserID, &p.Name, &p.Address, &p.Type, &p.PurchasePrice, &p.CurrentValue, &p.LoanID, &p.CreatedAt, &p.UpdatedAt)
}
