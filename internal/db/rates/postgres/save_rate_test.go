package postgresratesdb

import (
	"context"
	"itspay/internal/entity"
	"itspay/internal/utils/testkit"
	"time"
)

func (suite *DBTestSuite) TestSaveRate() {
	expectedRate := &entity.Rate{
		Ask:        testkit.MustDecimalFromString("1"),
		Bid:        testkit.MustDecimalFromString("2"),
		ReceivedAt: time.Now(),
	}

	err := suite.db.SaveRate(context.Background(), expectedRate)
	suite.Require().NoError(err)

	actualRate, err := suite.db.getLastRate(context.Background())
	suite.Require().NoError(err)
	testkit.RequireEqualCmp(suite.T(), expectedRate, actualRate, testkit.DecimalComparer)
}
