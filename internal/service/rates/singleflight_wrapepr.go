package ratesservice

import (
	"context"
	"itspay/internal/entity"
	"itspay/internal/service"

	"golang.org/x/sync/singleflight"
)

// SingleFlightService wraps Rates service to avoid unnecessary calls with the same result
// NOTE: this wrapper could have been in a separate package as with internal/rateprovider/metrics, both options possible
type SingleFlightService struct {
	s service.Rates
	g singleflight.Group
}

func NewSingleflightService(s service.Rates) *SingleFlightService {
	return &SingleFlightService{s: s}
}

func (s *SingleFlightService) GetRate(ctx context.Context) (*entity.Rate, error) {
	rate, err, _ := s.g.Do("GetRate", func() (any, error) {
		return s.s.GetRate(ctx)
	})

	return rate.(*entity.Rate), err
}
