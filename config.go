package mypggorm

import (
	"io"

	"github.com/go-faster/errors"
	"gopkg.in/yaml.v3"
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
