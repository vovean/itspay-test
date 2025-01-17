package api

import (
	"context"
	"fmt"
	"itspay/internal/config"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getDBConnectionURL(config *config.PostgresConfig) string {
	return fmt.Sprintf("postgresql://%s:%s@%s/%s?timezone=%s&sslmode=%s",
		config.User,
		config.Password,
		config.Addr,
		config.DB,
		"UTC",
		"disable",
	)
}

func setupPgxPool(ctx context.Context, c *config.PostgresConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(getDBConnectionURL(c))
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pool, nil
}
