package ratesservice

import (
	"context"
	"fmt"
	"itspay/internal/db"
	"itspay/internal/entity"
	"itspay/internal/rateprovider"
)

type Service struct {
	rateProvider rateprovider.RateProvider
	db           db.Rates
}

func New(rateProvider rateprovider.RateProvider, db db.Rates) *Service {
	return &Service{rateProvider: rateProvider, db: db}
}

func (s *Service) GetRate(ctx context.Context) (*entity.Rate, error) {
	rate, err := s.rateProvider.GetRate(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get rate: %w", err)
	}

	err = s.db.SaveRate(ctx, rate) // TODO: must it really reside on critical path? Discuss with business
	if err != nil {
		return nil, fmt.Errorf("cannot save rate: %w", err)
	}

	return rate, nil
}
