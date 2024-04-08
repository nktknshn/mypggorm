package databasetest

import (
	"github.com/go-faster/errors"
	"github.com/nktknshn/mypggorm"
)

func SetupPostgres(rootPassword string) (*DockerDatabase, error) {
	dt := NewDockerDatabase(
		DatabaseTestDockerConfig{
			DockerPassword: rootPassword,
			Timezone:       defaultDockerConfig.Timezone,
			Repository:     defaultDockerConfig.Repository,
			Tag:            defaultDockerConfig.Tag,
			RestartPolicy:  defaultDockerConfig.RestartPolicy,
			ResourceExpire: defaultDockerConfig.ResourceExpire,
		},
	)

	if err := dt.RunPostgresDocker(); err != nil {
		return nil, errors.Wrap(err, "failed to run postgres docker")
	}

	return dt, nil
}

func SetupDatabase(userCfg mypggorm.DatabaseConnectionConfig) (*DockerDatabase, error) {
	dt := NewDockerDatabase(
		DatabaseTestDockerConfig{},
	)

	if err := dt.RunPostgresDocker(); err != nil {
		return nil, err
	}

	if err := dt.CreateUserAndDatabase(userCfg); err != nil {
		return nil, err
	}

	return dt, nil
}
