package api

import (
	"context"
	"flag"
	"itspay/internal/config"
	"itspay/internal/utils/configkit"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

func applyFlagsToConfig(c *config.Config) {
	// It is only required to have flags for postgres so I define them here manually
	// If I need a proper solution to have flags for every field in the config I would define another backend
	// like for env variables which
	flag.StringVar(&c.Postgres.Addr, "postgres.addr", c.Postgres.Addr, "PostgreSQL address (host:port)")
	flag.StringVar(&c.Postgres.DB, "postgres.db", c.Postgres.DB, "PostgreSQL database")
	flag.StringVar(&c.Postgres.User, "postgres.user", c.Postgres.User, "PostgreSQL user")
	flag.StringVar(&c.Postgres.Password, "postgres.password", c.Postgres.Password, "PostgreSQL password")

	flag.Parse()
}

func loadConfig(ctx context.Context) (*config.Config, error) {
	loader := confita.NewLoader(
		file.NewBackend("config/rates_api.yml"),
		configkit.NewNestedEnvBackend(),
	)

	var c config.Config
	if err := loader.Load(ctx, &c); err != nil {
		return nil, err
	}

	applyFlagsToConfig(&c)

	return &c, nil
}
