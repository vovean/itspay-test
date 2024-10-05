package rateprovider

import (
	"context"
	"itspay/internal/entity"
)

//go:generate moq -fmt goimports -out mock/rate_provider.go -pkg mockrateprovider . RateProvider
type RateProvider interface {
	GetRate(ctx context.Context) (*entity.Rate, error)
}
