package mypggorm

import (
	"strings"

	"gorm.io/gorm"

	pg "gorm.io/driver/postgres"
)

type DatabaseConnection struct {
	connection *gorm.DB
}

func NewDatabaseConnection(connection *gorm.DB) *DatabaseConnection {
	return &DatabaseConnection{connection: connection}
}

func (dc *DatabaseConnection) Connection() *gorm.DB {
	return dc.connection
}

func (dc *DatabaseConnection) JustExec(sql string) error {
	return dc.connection.Exec(sql).Error
}

func (dc *DatabaseConnection) Close() error {
	sqlDB, err := dc.connection.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

type DatabaseConnectionConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
	Timezone string `yaml:"timezone"`
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

func (cfg DatabaseConnectionConfig) ToPostgresConfig() pg.Config {
	return pg.Config{
		DSN: cfg.DSN(),
	}
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
