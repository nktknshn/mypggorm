package postgres

// pg "gorm.io/driver/postgres"

import (
	"fmt"

	"github.com/nktknshn/mypggorm"

	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresInstanceConnection struct {
	mypggorm.DatabaseConnection
}

type PostgresInstance struct {
	rootConnectionConfig mypggorm.DatabaseConnectionConfig
}

func NewPostgresInstance(rootConnectionConfig mypggorm.DatabaseConnectionConfig) *PostgresInstance {
	return &PostgresInstance{
		rootConnectionConfig: rootConnectionConfig,
	}
}

func (p *PostgresInstance) DropDatabase(cfg *mypggorm.DatabaseConnectionConfig, dropUser bool) error {

	if cfg.Dbname == "" || cfg.Dbname == "postgres" {
		return fmt.Errorf("cannot drop database %s", cfg.Dbname)
	}

	conn, err := p.ConnectAsRoot("postgres")

	if err != nil {
		return err
	}

	err = conn.JustExec(fmt.Sprintf("DROP DATABASE %s;", cfg.Dbname))

	if err != nil {
		return err
	}

	// delete user
	if dropUser {
		err = conn.JustExec(fmt.Sprintf("DROP USER %s;", cfg.User))

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresInstance) CheckRootConnection() error {
	conn, err := p.ConnectAsRoot("postgres")

	if err != nil {
		return err
	}

	return conn.Close()
}

func (p *PostgresInstance) CheckConnection(cfg mypggorm.DatabaseConnectionConfig) error {
	conn, err := p.ConnectDatabase(cfg.PostgresConfig())

	if err != nil {
		return err
	}

	return conn.Close()
}

func (p *PostgresInstance) CreateUserAndDatabase(cfg mypggorm.DatabaseConnectionConfig, createUser bool) error {

	conn, err := p.ConnectAsRoot("postgres")

	if err != nil {
		return err
	}

	err = conn.JustExec("CREATE DATABASE " + cfg.Dbname + ";")

	if err != nil {
		return err
	}

	if createUser {
		err = conn.JustExec("CREATE USER " + cfg.User + " WITH PASSWORD '" + cfg.Password + "';")

		if err != nil {
			return err
		}
	}

	err = conn.JustExec("GRANT ALL ON DATABASE " + cfg.Dbname + " TO " + cfg.User + ";")

	if err != nil {
		return err
	}

	if err := conn.Close(); err != nil {
		return err
	}

	// GRANT ON SCHEMA
	conn, err = p.ConnectAsRoot(cfg.Dbname)

	if err != nil {
		return err
	}

	err = conn.JustExec("GRANT ALL ON SCHEMA public TO " + cfg.User + ";")

	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresInstance) CleanDatabase(cfg mypggorm.DatabaseConnectionConfig) error {
	conn, err := p.ConnectAsRoot(cfg.Dbname)

	/*
			-- Recreate the schema
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;

		-- Restore default permissions
		GRANT ALL ON SCHEMA public TO postgres;
		GRANT ALL ON SCHEMA public TO public;
	*/
	if err != nil {
		return err
	}

	if err := conn.JustExec(fmt.Sprintf("DROP SCHEMA public CASCADE;")); err != nil {
		return err
	}

	if err := conn.JustExec(fmt.Sprintf("CREATE SCHEMA public;")); err != nil {
		return err
	}

	if err := conn.JustExec(fmt.Sprintf("GRANT ALL ON SCHEMA public TO postgres;")); err != nil {
		return err
	}

	if err := conn.JustExec(fmt.Sprintf("GRANT ALL ON SCHEMA public TO %s;", cfg.User)); err != nil {
		return err
	}

	if err := conn.JustExec(fmt.Sprintf("GRANT ALL ON SCHEMA public TO public;")); err != nil {
		return err
	}

	return conn.Close()
}

func (p *PostgresInstance) ConnectAsRoot(dbname string) (*PostgresInstanceConnection, error) {

	return p.ConnectDatabase(pg.Config{
		DSN: p.rootConnectionConfig.WithDbname(dbname).DSN(),
	})
}

func (p *PostgresInstance) ConnectDatabase(pgConfig pg.Config) (*PostgresInstanceConnection, error) {

	conn, err := gorm.Open(pg.New(pgConfig), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &PostgresInstanceConnection{
		DatabaseConnection: *mypggorm.NewDatabaseConnection(conn),
	}, nil
}

func (p *PostgresInstance) BackupDatabase(cfg mypggorm.DatabaseConnectionConfig, backupPath string) error {
	return nil
}

func (p *PostgresInstance) TimeZone() (string, error) {

	conn, err := p.ConnectAsRoot("postgres")

	if err != nil {
		return "", err
	}

	var tz string
	conn.DatabaseConnection.Connection().Raw("SELECT current_setting('TIMEZONE')").Scan(&tz)

	err = conn.Close()

	if err != nil {
		return "", err
	}

	return tz, nil

}
