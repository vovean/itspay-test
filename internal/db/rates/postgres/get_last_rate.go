package postgresratesdb

import (
	"context"
	_ "embed"
	"itspay/internal/entity"
)

//go:embed sql/get_last_rate.sql
var getLastRateQuery string

// getLastRate used only in tests.
func (db *DB) getLastRate(ctx context.Context) (*entity.Rate, error) {
	rate := &entity.Rate{}

	err := db.pool.QueryRow(ctx, getLastRateQuery).Scan(&rate.Ask, &rate.Bid, &rate.ReceivedAt)
	if err != nil {
		return nil, err
	}

	return rate, nil
}
