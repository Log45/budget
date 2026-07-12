package db

import "github.com/jackc/pgx/v5/pgxpool"

type TransactionRepository struct {
	pool *pgxpool.Pool
}
