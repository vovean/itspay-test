package postgresratesdb

import (
	"context"
	_ "embed"
	"itspay/migrations"
	"os"
	"testing"

	"github.com/antelman107/net-wait-go/wait"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

var (
	connString = "postgresql://postgres:postgres@localhost:5433/postgres?sslmode=disable"
	pool       *pgxpool.Pool
)

func TestMain(m *testing.M) {
	if !wait.New().Do([]string{"localhost:5433"}) {
		panic("can't await db")
	}

	var err error

	pool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.Exit(code)
}

type baseDBTestSuite struct {
	suite.Suite
	db       *DB
	migrator *migrate.Migrate
}

func (suite *baseDBTestSuite) SetupSuite() {
	suite.db = New(pool)

	ioFSSourceInstance, err := iofs.New(migrations.FS, ".")
	suite.Require().NoError(err)

	suite.migrator, err = migrate.NewWithSourceInstance(
		"iofs",
		ioFSSourceInstance,
		connString,
	)
	suite.Require().NoError(err)
}

func (suite *baseDBTestSuite) SetupTest() {
	suite.Require().NoError(suite.migrator.Up())
}

func (suite *baseDBTestSuite) TearDownTest() {
	suite.Require().NoError(suite.migrator.Down())
}
