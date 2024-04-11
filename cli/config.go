package cli

import (
	"fmt"
	"os"

	"github.com/go-faster/errors"
	"github.com/nktknshn/mypggorm"
)

type DatabaseConfigProvider struct {
	ConfigFile         func() string
	RootPasswordEnvKey string
}

func NewDatabaseConfigProvider(configFile func() string) *DatabaseConfigProvider {
	return &DatabaseConfigProvider{
		ConfigFile:         configFile,
		RootPasswordEnvKey: "POSTGRES_PASSWORD",
	}
}

func (d *DatabaseConfigProvider) GetUserConfig() (mypggorm.DatabaseConnectionConfig, error) {
	r, err := os.Open(d.ConfigFile())
	if err != nil {
		return mypggorm.DatabaseConnectionConfig{}, errors.Wrap(err, "failed to read config file")
	}
	cfg, err := mypggorm.ParseYAMLConfig(r)
	if err != nil {
		return mypggorm.DatabaseConnectionConfig{}, errors.Wrap(err, "failed to parse config file")
	}
	return cfg, nil
}

func (d *DatabaseConfigProvider) GetRootConfig() (mypggorm.DatabaseConnectionConfig, error) {
	pw := os.Getenv(d.RootPasswordEnvKey)
	if pw == "" {
		return mypggorm.DatabaseConnectionConfig{}, fmt.Errorf("POSTGRES_PASSWORD is not set")
	}
	cfg, err := d.GetUserConfig()
	if err != nil {
		return mypggorm.DatabaseConnectionConfig{}, errors.Wrap(err, "failed to get user config")
	}
	return cfg.WithUser("postgres").WithPassword(pw).WithDbname("postgres"), nil
}
