package db

import "github.com/jackc/pgx/v5/pgxpool"

type LoanRepository struct {
	pool *pgxpool.Pool
}
