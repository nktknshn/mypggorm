package main

import (
	"fmt"
	"os"

	"github.com/go-faster/errors"
	"github.com/nktknshn/mypggorm"
	dbcli "github.com/nktknshn/mypggorm/cli"
	"github.com/spf13/cobra"
)

var (
	flagDBConfig string = "database-config.yaml"
)

type DatabaseConfig struct{}

func (d *DatabaseConfig) GetUserConfig() (mypggorm.DatabaseConnectionConfig, error) {
	r, err := os.Open(flagDBConfig)
	if err != nil {
		return mypggorm.DatabaseConnectionConfig{}, errors.Wrap(err, "failed to read config file")
	}
	cfg, err := mypggorm.ParseYAMLConfig(r)
	if err != nil {
		return mypggorm.DatabaseConnectionConfig{}, errors.Wrap(err, "failed to parse config file")
	}
	return cfg, nil
}

func (d *DatabaseConfig) GetRootConfig() (mypggorm.DatabaseConnectionConfig, error) {

	pw := os.Getenv("POSTGRES_PASSWORD")
	if pw == "" {
		return mypggorm.DatabaseConnectionConfig{}, fmt.Errorf("POSTGRES_PASSWORD is not set")
	}
	cfg, err := d.GetUserConfig()
	if err != nil {
		return mypggorm.DatabaseConnectionConfig{}, errors.Wrap(err, "failed to get user config")
	}
	return cfg.SetPassword(pw).SetDbname("postgres"), nil
}

func init() {
	commandDB := dbcli.CreateCommand(&DatabaseConfig{})
	commandDB.PersistentFlags().StringVarP(&flagDBConfig, "config", "c", flagDBConfig, "database config file")
	cmd.AddCommand(commandDB)
}

var cmd = &cobra.Command{
	Use:  "cli",
	Args: cobra.MinimumNArgs(1),
}

func main() {
	err := cmd.Execute()

	if err != nil {
		panic(err)
	}
}
