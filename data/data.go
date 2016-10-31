package data

import "github.com/jackc/pgx"

var (
	pool *pgx.ConnPool
)

// Pool configures the PostgreSQL pool
func Pool(use *pgx.ConnPool) {
	pool = use
}
