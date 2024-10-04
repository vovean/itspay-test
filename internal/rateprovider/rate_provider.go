package rateprovider

import (
	"context"
	"itspay/internal/entity"
)

type RateProvider interface {
	GetRate(ctx context.Context) (*entity.Rate, error)
}
