package ratesservice

import (
	"context"
	"errors"
	"itspay/internal/entity"
	mockservice "itspay/internal/service/mock"
	"itspay/internal/utils/testkit"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type SingleFlightWrapperTestSuite struct {
	suite.Suite
	service *mockservice.RatesMock
	wrapper *SingleFlightService
}

func (suite *SingleFlightWrapperTestSuite) SetupTest() {
	suite.service = &mockservice.RatesMock{}
	suite.wrapper = NewSingleflightService(suite.service)
}

func (suite *SingleFlightWrapperTestSuite) TestSuccess() {
	expectedRate := &entity.Rate{
		Ask:        testkit.MustDecimalFromString("1"),
		Bid:        testkit.MustDecimalFromString("2"),
		ReceivedAt: time.Now(),
	}
	suite.service.GetRateFunc = func(ctx context.Context) (*entity.Rate, error) {
		return expectedRate, nil
	}

	actualRate, err := suite.wrapper.GetRate(context.Background())
	suite.Require().NoError(err)
	testkit.RequireEqualCmp(suite.T(), expectedRate, actualRate, testkit.DecimalComparer)
}

func (suite *SingleFlightWrapperTestSuite) TestError() {
	expectedError := errors.New("some error")
	suite.service.GetRateFunc = func(ctx context.Context) (*entity.Rate, error) {
		return nil, expectedError
	}

	_, err := suite.wrapper.GetRate(context.Background())
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, expectedError)
}

func (suite *SingleFlightWrapperTestSuite) TestSingleFlight() {
	expectedRate := &entity.Rate{
		Ask:        testkit.MustDecimalFromString("1"),
		Bid:        testkit.MustDecimalFromString("2"),
		ReceivedAt: time.Now(),
	}
	suite.service.GetRateFunc = func(ctx context.Context) (*entity.Rate, error) {
		time.Sleep(time.Second)

		return expectedRate, nil
	}

	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			actualRate, err := suite.wrapper.GetRate(context.Background())
			suite.NoError(err)
			testkit.AssertEqualCmp(suite.T(), expectedRate, actualRate, testkit.DecimalComparer)
		}()
	}

	wg.Wait()

	suite.Require().Len(suite.service.GetRateCalls(), 1)
}

func TestSingleFlightWrapperTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(SingleFlightWrapperTestSuite))
}
