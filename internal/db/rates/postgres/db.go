package postgresratesdb

import (
	ratesdb "itspay/internal/db/rates"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ ratesdb.DB = &DB{}

type DB struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *DB {
	return &DB{pool: pool}
}
