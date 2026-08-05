package config

import "time"

type (
	APIConfig struct {
		ShutdownTimeout time.Duration  `yaml:"shutdown_timeout"`
		LoggerConfig    LoggerConfig   `yaml:"logger"`
		PostgresConfig  PostgresConfig `yaml:"postgres"`
		ServerConfig    ServerConfig   `yaml:"server"`
		JWTConfig       JWTConfig      `yaml:"jwt"`
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

	JWTConfig struct {
		PrivateKeyPath string        `yaml:"private_key_path" env:"JWT_PRIVATE_KEY_PATH"`
		PublicKeyPath  string        `yaml:"public_key_path" env:"JWT_PUBLIC_KEY_PATH"`
		Issuer         string        `yaml:"issuer"`
		Audience       string        `yaml:"audience"`
		KeyID          string        `yaml:"key_id"`
		TTL            time.Duration `yaml:"ttl"`
	}
)
