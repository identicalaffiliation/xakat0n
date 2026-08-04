package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type MigratorConfig struct {
	LoggerConfig LoggerConfig `yaml:"logger"`
	DatabaseUrl  string       `env:"DB_URL"`
}

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func LoadConfig(configPath string) (*MigratorConfig, error) {
	cfg := new(MigratorConfig)
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return cfg, nil
}
