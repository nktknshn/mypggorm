package databasetest

import (
	"errors"

	"github.com/nktknshn/mypggorm"
	dockertest "github.com/ory/dockertest/v3"
	docker "github.com/ory/dockertest/v3/docker"
	pg "gorm.io/driver/postgres"

	"gorm.io/gorm"
)

type DatabaseTestDockerConfig struct {
	Repository     string
	Tag            string
	DockerPassword string
	Timezone       string
	RestartPolicy  docker.RestartPolicy
	ResourceExpire uint
}

var defaultDockerConfig = DatabaseTestDockerConfig{
	Repository:     "postgres",
	Tag:            "latest",
	DockerPassword: "who-cares-about-the-test-password",
	RestartPolicy:  docker.RestartPolicy{Name: "no"},
	ResourceExpire: 120,
	// Timezone:       "UTC",
	Timezone: "Europe/Moscow",
}

type RunningPostgres struct {
	Container *dockertest.Resource
}

func (rp *RunningPostgres) GetPort(name string) string {
	return rp.Container.GetPort(name)
}

type DockerDatabase struct {
	dockerConfig DatabaseTestDockerConfig

	runningDocker *RunningPostgres
	pool          *dockertest.Pool
}

func NewDockerDatabase(dockerConfig DatabaseTestDockerConfig) *DockerDatabase {
	return &DockerDatabase{
		dockerConfig: dockerConfig,
	}
}

func (dt *DockerDatabase) IsRunning() bool {
	return dt.runningDocker != nil
}

func (dt *DockerDatabase) GetRunningPort() (string, error) {
	if dt.runningDocker == nil {
		return "", errors.New("postgres docker is not running")
	}

	return dt.runningDocker.GetPort("5432/tcp"), nil
}

func (dt *DockerDatabase) GetPool() *dockertest.Pool {
	return dt.pool
}

func (dt *DockerDatabase) RunPostgresDocker() error {

	pool, err := dockertest.NewPool("")

	dt.pool = pool

	if err != nil {
		return err
	}

	if err = pool.Client.Ping(); err != nil {
		return err
	}

	if dt.dockerConfig.DockerPassword == "" {
		dt.dockerConfig.DockerPassword = defaultDockerConfig.DockerPassword
	}

	if dt.dockerConfig.Repository == "" {
		dt.dockerConfig.Repository = defaultDockerConfig.Repository
	}

	if dt.dockerConfig.Tag == "" {
		dt.dockerConfig.Tag = defaultDockerConfig.Tag
	}

	if dt.dockerConfig.RestartPolicy == (docker.RestartPolicy{}) {
		dt.dockerConfig.RestartPolicy = defaultDockerConfig.RestartPolicy
	}

	if dt.dockerConfig.ResourceExpire == 0 {
		dt.dockerConfig.ResourceExpire = defaultDockerConfig.ResourceExpire
	}

	if dt.dockerConfig.Timezone == "" {
		dt.dockerConfig.Timezone = defaultDockerConfig.Timezone
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: dt.dockerConfig.Repository,
		Tag:        dt.dockerConfig.Tag,
		Env: []string{
			"POSTGRES_PASSWORD=" + dt.dockerConfig.DockerPassword,
			"TZ=" + dt.dockerConfig.Timezone,
			"PGTZ=" + dt.dockerConfig.Timezone,
		},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = dt.dockerConfig.RestartPolicy
	})

	if err != nil {
		return err
	}

	if err = resource.Expire(dt.dockerConfig.ResourceExpire); err != nil {
		return err
	}

	dt.runningDocker = &RunningPostgres{
		Container: resource,
	}

	return nil
}

func (dt *DockerDatabase) ConnectDatabaseConfig(config mypggorm.DatabaseConnectionConfig) (*mypggorm.DatabaseConnection, error) {
	return dt.connectDatabase(pg.Config{
		DSN: config.SetPort(dt.runningDocker.Container.GetPort("5432/tcp")).DSN(),
	})
}

func (dt *DockerDatabase) connectDatabase(pgConfig pg.Config) (*mypggorm.DatabaseConnection, error) {

	var dbCon *gorm.DB

	if err := dt.pool.Retry(func() (err error) {
		dbCon, err = gorm.Open(
			pg.New(pgConfig),
			&gorm.Config{},
		)
		return
	}); err != nil {
		return nil, err
	}

	return mypggorm.NewDatabaseConnection(dbCon), nil
}

func (dt *DockerDatabase) ConnectDatabaseAsRoot(dbName string) (*mypggorm.DatabaseConnection, error) {

	pgConfig := pg.Config{
		DSN: mypggorm.DefaultDatabaseConnectionConfig.SetDbname(dbName).SetPassword(dt.dockerConfig.DockerPassword).
			SetPort(dt.runningDocker.Container.GetPort("5432/tcp")).SetUser("postgres").DSN(),
	}

	return dt.connectDatabase(pgConfig)
}

func (dt *DockerDatabase) StopPostgresDocker() error {

	if dt.runningDocker == nil {
		return nil
	}

	d := dt.runningDocker

	dt.runningDocker = nil

	return dt.pool.Purge(d.Container)
}
