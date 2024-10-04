package db

import (
	"context"
	"itspay/internal/entity"
)

type Rates interface {
	SaveRate(ctx context.Context, rate *entity.Rate) error
}
