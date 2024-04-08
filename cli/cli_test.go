package cli_test

import (
	"testing"

	"github.com/go-faster/errors"
	"github.com/nktknshn/mypggorm"
	"github.com/nktknshn/mypggorm/cli"
	"github.com/nktknshn/mypggorm/databasetest"
	"github.com/stretchr/testify/suite"
)

type MyTestSuite struct {
	suite.Suite

	DT  *databasetest.DockerDatabase
	cfg TestCfg
}

type TestCfg struct {
	RootPassword string
	mypggorm.DatabaseConnectionConfig
}

var testCfg = TestCfg{
	RootPassword: "testrootpassword",
	DatabaseConnectionConfig: mypggorm.DatabaseConnectionConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpassword",
		Dbname:   "testdb",
		Sslmode:  "disable",
		Timezone: "UTC",
	},
}

func (t TestCfg) GetUserConfig() (mypggorm.DatabaseConnectionConfig, error) {
	return t.DatabaseConnectionConfig, nil
}

func (t TestCfg) GetRootConfig() (mypggorm.DatabaseConnectionConfig, error) {
	return t.DatabaseConnectionConfig.WithPassword(t.RootPassword).WithUser("postgres"), nil
}

func (suite *MyTestSuite) SetupSuite() {
	suite.T().Log("Setting up docker")
	suite.cfg = testCfg
	dt, err := databasetest.SetupPostgres(suite.cfg.RootPassword)
	if err != nil {
		suite.T().Fatal(errors.Wrap(err, "failed to setup database"))
	}
	suite.DT = dt
	port, err := suite.DT.GetRunningPort()
	suite.T().Logf("Running port: %s", port)
	if err != nil {
		suite.T().Fatal(errors.Wrap(err, "failed to get running port"))
	}
	suite.cfg.DatabaseConnectionConfig = suite.cfg.WithPort(port)

	// wait for postgres to start
	if _, err = suite.DT.ConnectDatabaseAsRoot("postgres"); err != nil {
		suite.T().Fatal(errors.Wrap(err, "failed to connect to database"))
	}
}

func (suite *MyTestSuite) TearDownSuite() {
	suite.T().Log("Stopping docker")
	if err := suite.DT.StopPostgresDocker(); err != nil {
		suite.FailNow(err.Error())
	}
}

func (suite *MyTestSuite) TearDownTest() {
	suite.T().Log("Cleaning database")
}

type TestCliSuite struct {
	MyTestSuite
}

func (s *TestCliSuite) TestCreate() {
	cmd := cli.DatabaseCommand(s.cfg)
	cmd.SetArgs([]string{"create"})
	if err := cmd.Execute(); err != nil {
		s.FailNow(err.Error())
	}
}

func TestCli(t *testing.T) {
	suite.Run(t, new(TestCliSuite))
}
