package ratesapi

import (
	"context"
	"errors"
	"itspay/internal/api/rates/ratespb"
	"itspay/internal/entity"
	"itspay/internal/utils/testkit"
	"time"

	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (suite *ServerTestSuite) TestGetRate_Success() {
	rate := &entity.Rate{
		Ask:        testkit.MustDecimalFromString("1"),
		Bid:        testkit.MustDecimalFromString("2"),
		ReceivedAt: time.Now(),
	}
	suite.service.GetRateFunc = func(ctx context.Context) (*entity.Rate, error) {
		return rate, nil
	}

	expectedRate := &ratespb.GetRateResponse{
		Ask:        "1",
		Bid:        "2",
		ReceivedAt: timestamppb.New(rate.ReceivedAt),
	}

	actualRate, err := suite.server.GetRate(context.Background(), &ratespb.GetRateRequest{})
	suite.Require().NoError(err)
	testkit.RequireEqualCmp(suite.T(), expectedRate, actualRate, protocmp.Transform())
}

func (suite *ServerTestSuite) TestGetRate_ServiceError() {
	suite.service.GetRateFunc = func(ctx context.Context) (*entity.Rate, error) {
		return nil, errors.New("service error")
	}

	_, err := suite.server.GetRate(context.Background(), &ratespb.GetRateRequest{})
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, errInternal)
}
