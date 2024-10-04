package service

import (
	"context"
	"itspay/internal/entity"
)

type Rates interface {
	GetRate(ctx context.Context) (*entity.Rate, error)
}
