package ratesservice

import (
	"context"
	"errors"
	mockratesdb "itspay/internal/db/rates/mock"
	"itspay/internal/entity"
	mockrateprovider "itspay/internal/rateprovider/mock"
	"itspay/internal/utils/testkit"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	rateProvider *mockrateprovider.RateProviderMock
	db           *mockratesdb.DBMock
	service      *Service
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.rateProvider = &mockrateprovider.RateProviderMock{}
	suite.db = &mockratesdb.DBMock{}
	suite.service = New(suite.rateProvider, suite.db)
}

func (suite *ServiceTestSuite) TestGetRate_Success() {
	expectedRate := &entity.Rate{
		Ask:        testkit.MustDecimalFromString("1"),
		Bid:        testkit.MustDecimalFromString("2"),
		ReceivedAt: time.Now(),
	}
	suite.rateProvider.GetRateFunc = func(ctx context.Context) (*entity.Rate, error) {
		return expectedRate, nil
	}

	suite.db.SaveRateFunc = func(ctx context.Context, rate *entity.Rate) error {
		testkit.AssertEqualCmp(suite.T(), expectedRate, rate, testkit.DecimalComparer)

		return nil
	}

	actualRate, err := suite.service.GetRate(context.Background())
	suite.Require().NoError(err)
	testkit.RequireEqualCmp(suite.T(), expectedRate, actualRate, testkit.DecimalComparer)
}

func (suite *ServiceTestSuite) TestGetRate_RateProviderError() {
	expectedError := errors.New("some error")
	suite.rateProvider.GetRateFunc = func(ctx context.Context) (*entity.Rate, error) {
		return nil, expectedError
	}

	_, err := suite.service.GetRate(context.Background())
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, expectedError)
}

func (suite *ServiceTestSuite) TestGetRate_DBError() {
	suite.rateProvider.GetRateFunc = func(ctx context.Context) (*entity.Rate, error) {
		return &entity.Rate{}, nil
	}

	expectedError := errors.New("some error")
	suite.db.SaveRateFunc = func(ctx context.Context, rate *entity.Rate) error {
		return expectedError
	}

	_, err := suite.service.GetRate(context.Background())
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, expectedError)
}

func TestServiceTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ServiceTestSuite))
}
