package postgresratesdb

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	baseDBTestSuite
}

func TestDBTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(DBTestSuite))
}
