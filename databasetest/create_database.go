package databasetest

import (
	"fmt"

	"github.com/nktknshn/mypggorm"
)

func (dt *DockerDatabase) CreateUserAndDatabase(cfg mypggorm.DatabaseConnectionConfig) error {
	conn, err := dt.ConnectDatabaseAsRoot("postgres")

	if err != nil {
		return err
	}

	err = conn.JustExec(fmt.Sprintf("CREATE DATABASE %s;", cfg.Dbname))

	if err != nil {
		return err
	}

	err = conn.JustExec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s';", cfg.User, cfg.Password))

	if err != nil {
		return nil
	}

	err = conn.JustExec(fmt.Sprintf("GRANT ALL ON DATABASE %s TO %s;", cfg.Dbname, cfg.User))

	if err != nil {
		return err
	}

	if err := conn.Close(); err != nil {
		return err
	}

	// GRANT ON SCHEMA
	conn, err = dt.ConnectDatabaseAsRoot(cfg.Dbname)

	if err != nil {
		return err
	}

	err = conn.JustExec(fmt.Sprintf("GRANT ALL ON SCHEMA public TO %s;", cfg.User))

	if err != nil {
		return nil
	}

	return conn.Close()

}

func (dt *DockerDatabase) CleanDatabase(cfg mypggorm.DatabaseConnectionConfig) error {
	conn, err := dt.ConnectDatabaseAsRoot(cfg.Dbname)

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
