package server

import (
	env "github.com/caarlos0/env/v6"
)

// Config contains the runtime config options for this service
type Config struct {
	Port           int    `env:"PORT" envDefault:"8443"`
	LabelsB64      string `env:"K8S_LABELS" envDefault:""`
	AnnotationsB64 string `env:"K8S_ANNOTATIONS" envDefault:""`
}

// ConfigFromEnv returns a Config instance from env vars
func ConfigFromEnv() (*Config, error) {
	cfg := &Config{}
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.Parse(cfg, opts); err != nil {
		return nil, err
	}
	return cfg, nil
}
