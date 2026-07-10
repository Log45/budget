package db

import "github.com/jackc/pgx/v5/pgxpool"

type BudgetRepository struct {
	pool *pgxpool.Pool
}
