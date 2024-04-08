package databasetest

import "github.com/nktknshn/mypggorm"

func SetupDatabase(cfg mypggorm.DatabaseConnectionConfig) (*DockerDatabase, error) {
	dt := NewDockerDatabase(
		DatabaseTestDockerConfig{},
	)

	if err := dt.RunPostgresDocker(); err != nil {
		return nil, err
	}

	if err := dt.CreateUserAndDatabase(cfg); err != nil {
		return nil, err
	}

	return dt, nil
}
