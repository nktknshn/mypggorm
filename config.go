package mypggorm

import (
	"io"
	"strings"
	"time"

	"github.com/go-faster/errors"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ParseYAMLConfig(r io.Reader) (DatabaseConnectionConfig, error) {
	cfg := DatabaseConnectionConfig{}

	err := yaml.NewDecoder(r).Decode(&cfg)

	if err != nil {
		return DatabaseConnectionConfig{}, errors.Wrap(err, "failed to decode yaml")
	}

	if cfg.Host == "" {
		return DatabaseConnectionConfig{}, errors.New("Invalid config file: host is empty")
	}

	return cfg, nil
}

type DatabaseConnectionConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
	Timezone string `yaml:"TimeZone"`

	NowFunc func() time.Time `yaml:"-"`
	DoLog   bool
}

func (cfg DatabaseConnectionConfig) Connect() (*gorm.DB, error) {
	gcfg := cfg.GormConfig()
	return gorm.Open(
		postgres.New(cfg.PostgresConfig()),
		gcfg,
	)
}

func (cfg DatabaseConnectionConfig) PostgresConfig() postgres.Config {
	return postgres.Config{
		DSN: cfg.DSN(),
	}
}

// ToString
func (cfg DatabaseConnectionConfig) String() string {
	return "host=" + cfg.Host +
		" user=" + cfg.User +
		" dbname=" + cfg.Dbname +
		" port=" + cfg.Port +
		" sslmode=" + cfg.Sslmode +
		" TimeZone=" + cfg.Timezone
}

func (cfg DatabaseConnectionConfig) GormConfig() *gorm.Config {
	gormConfig := gorm.Config{}

	if cfg.NowFunc != nil {
		gormConfig.NowFunc = cfg.NowFunc
	}

	if cfg.DoLog {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	return &gormConfig
}

func ParseDSN(dsn string) DatabaseConnectionConfig {
	cfg := DatabaseConnectionConfig{}

	for _, part := range strings.Split(dsn, " ") {
		parts := strings.Split(part, "=")

		if len(parts) != 2 {
			continue
		}

		switch parts[0] {
		case "host":
			cfg.Host = parts[1]
		case "password":
			cfg.Password = parts[1]
		case "user":
			cfg.User = parts[1]
		case "dbname":
			cfg.Dbname = parts[1]
		case "port":
			cfg.Port = parts[1]
		case "sslmode":
			cfg.Sslmode = parts[1]
		case "TimeZone":
			cfg.Timezone = parts[1]
		}
	}

	return cfg
}

func (cfg DatabaseConnectionConfig) DSN() string {
	return "host=" + cfg.Host +
		" user=" + cfg.User +
		" password=" + cfg.Password +
		" dbname=" + cfg.Dbname +
		" port=" + cfg.Port +
		" sslmode=" + cfg.Sslmode +
		" TimeZone=" + cfg.Timezone
}

func (cfg DatabaseConnectionConfig) WithPort(port string) DatabaseConnectionConfig {
	cfg.Port = port
	return cfg
}

func (cfg DatabaseConnectionConfig) WithDbname(dbname string) DatabaseConnectionConfig {
	cfg.Dbname = dbname
	return cfg
}

func (cfg DatabaseConnectionConfig) WithPassword(password string) DatabaseConnectionConfig {
	cfg.Password = password
	return cfg
}

func (cfg DatabaseConnectionConfig) WithUser(user string) DatabaseConnectionConfig {
	cfg.User = user
	return cfg
}

var DefaultDatabaseConnectionConfig = DatabaseConnectionConfig{
	Host:     "localhost",
	User:     "postgres",
	Port:     "5432",
	Password: "postgres",
	Dbname:   "postgres",
	Sslmode:  "disable",
	Timezone: "UTC",
}
