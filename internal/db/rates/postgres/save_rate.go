package postgresratesdb

import (
	"context"
	_ "embed"
	"itspay/internal/entity"
)

//go:embed sql/save_rate.sql
var saveRateQuery string

func (db *DB) SaveRate(ctx context.Context, rate *entity.Rate) error {
	_, err := db.pool.Exec(ctx, saveRateQuery, rate.Ask, rate.Bid, rate.ReceivedAt)

	return err
}
