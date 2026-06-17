package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabasePath string `envconfig:"DATABASE_PATH" default:"./main.db"`
	DatabaseType string `envconfig:"DATABASE_TYPE" default:"sqlite3"`
}

func Load() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	return cfg, err
}
