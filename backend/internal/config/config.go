package config

import "time"

type (
	APIConfig struct {
		ShutdownTimeout time.Duration  `yaml:"shutdown_timeout"`
		LoggerConfig    LoggerConfig   `yaml:"logger"`
		PostgresConfig  PostgresConfig `yaml:"postgres"`
		ServerConfig    ServerConfig   `yaml:"server"`
	}

	LoggerConfig struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	}

	PostgresConfig struct {
		DbUrl       string        `env:"DB_URL"`
		MaxConns    int32         `yaml:"max_conns"`
		MaxLifetime time.Duration `yaml:"max_lifetime"`
	}

	ServerConfig struct {
		Host         string        `yaml:"host"`
		Port         int           `yaml:"port"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
		IddleTimeout time.Duration `yaml:"iddle_timeout"`
	}
)
