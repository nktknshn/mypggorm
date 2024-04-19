package cli

import (
	"fmt"
	"os"

	"github.com/go-faster/errors"
	"github.com/nktknshn/mypggorm"
)

const RootPasswordEnvKey = "POSTGRES_PASSWORD"

type DatabaseConfigProviderFile struct {
	ConfigFilePath     func() string
	RootPasswordEnvKey string
}

// parses specified config file and returns DatabaseConnectionConfig. Implemenets DatabaseConfigGetter
func NewDatabaseConfigProviderFile(configFile func() string) *DatabaseConfigProviderFile {
	return &DatabaseConfigProviderFile{
		ConfigFilePath:     configFile,
		RootPasswordEnvKey: RootPasswordEnvKey,
	}
}

func (d *DatabaseConfigProviderFile) GetUserConfig() (mypggorm.DatabaseConnectionConfig, error) {
	r, err := os.Open(d.ConfigFilePath())
	if err != nil {
		return mypggorm.DatabaseConnectionConfig{}, errors.Wrap(err, "failed to read config file")
	}
	cfg, err := mypggorm.ParseYAMLConfig(r)
	if err != nil {
		return mypggorm.DatabaseConnectionConfig{}, errors.Wrap(err, "failed to parse config file")
	}
	return cfg, nil
}

func (d *DatabaseConfigProviderFile) GetRootConfig() (mypggorm.DatabaseConnectionConfig, error) {
	pw := os.Getenv(d.RootPasswordEnvKey)
	if pw == "" {
		return mypggorm.DatabaseConnectionConfig{}, fmt.Errorf("%s is not set", d.RootPasswordEnvKey)
	}
	cfg, err := d.GetUserConfig()
	if err != nil {
		return mypggorm.DatabaseConnectionConfig{}, errors.Wrap(err, "failed to get user config")
	}
	return cfg.WithUser("postgres").WithPassword(pw).WithDbname("postgres"), nil
}
