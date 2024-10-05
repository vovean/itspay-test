package ratesapi

import (
	mockservice "itspay/internal/service/mock"
	"itspay/internal/utils/testkit"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	service *mockservice.RatesMock
	server  *Server
}

func (suite *ServerTestSuite) SetupTest() {
	suite.service = &mockservice.RatesMock{}
	suite.server = NewServer(suite.service, testkit.NewOTELZapTestLogger(suite.T()))
}

func TestServerTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ServerTestSuite))
}
