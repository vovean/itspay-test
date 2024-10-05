package ratesdb

import (
	"context"
	"itspay/internal/entity"
)

//go:generate moq -fmt goimports -out mock/db.go -pkg mockratesdb . DB
type DB interface {
	SaveRate(ctx context.Context, rate *entity.Rate) error
}
