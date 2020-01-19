package server

import (
	"time"

	envconfig "praslar.com/gotasma/internal/pkg/env"
)

type (
	// Config hold HTTP server configurations
	Config struct {
		Address           string        `envconfig:"HTTP_ADDRESS"`
		Port              int           `envconfig:"HTTP_PORT"`
		ReadTimeout       time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"10s"`
		WriteTimeout      time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"10s"`
		ReadHeaderTimeout time.Duration `envconfig:"HTTP_READ_HEADER_TIMEOUT" default:"10s"`
		ShutdownTimeout   time.Duration `envconfig:"HTTP_SHUTDOWN_TIMEOUT" default:"10s"`
	}
)

func LoadConfigFromEnv() Config {
	var conf Config
	envconfig.Load(&conf)
	return conf
}
