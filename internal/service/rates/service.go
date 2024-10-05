package ratesservice

import (
	"context"
	"fmt"
	"itspay/internal/db/rates"
	"itspay/internal/entity"
	"itspay/internal/rateprovider"
	"itspay/internal/service"
)

var _ service.Rates = &Service{}

type Service struct {
	rateProvider rateprovider.RateProvider
	db           ratesdb.DB
}

func New(rateProvider rateprovider.RateProvider, db ratesdb.DB) *Service {
	return &Service{
		rateProvider: rateProvider,
		db:           db,
	}
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
