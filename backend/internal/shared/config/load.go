package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

func LoadConfig(configPath string) (*APIConfig, error) {
	cfg := new(APIConfig)
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return cfg, nil
}
