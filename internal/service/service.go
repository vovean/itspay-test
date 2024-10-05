package service

import (
	"context"
	"itspay/internal/entity"
)

//go:generate moq -fmt goimports -out ./mock/rates.go -pkg mockservice . Rates
type Rates interface {
	GetRate(ctx context.Context) (*entity.Rate, error)
}
