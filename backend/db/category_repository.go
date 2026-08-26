package db

import (
	"Log45/budget/backend/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrCategoryNotFound = errors.New("category not found")

type CategoryRepository struct{ pool *pgxpool.Pool }

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}
func (r *CategoryRepository) List(ctx context.Context, userID int64) ([]models.Category, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,user_id,name,transaction_type,user_id IS NULL,created_at FROM categories WHERE user_id IS NULL OR user_id=$1 ORDER BY user_id NULLS FIRST,name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.Category{}
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.TransactionType, &c.IsDefault, &c.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}
func (r *CategoryRepository) Create(ctx context.Context, c *models.Category) error {
	return r.pool.QueryRow(ctx, `INSERT INTO categories (user_id,name,transaction_type) VALUES ($1,$2,$3) RETURNING id,created_at`, c.UserID, c.Name, c.TransactionType).Scan(&c.ID, &c.CreatedAt)
}
func (r *CategoryRepository) Accessible(ctx context.Context, userID, id int64) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM categories WHERE id=$1 AND (user_id IS NULL OR user_id=$2))`, id, userID).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, ErrCategoryNotFound
	}
	return exists, err
}
