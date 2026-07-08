package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig `yaml:"server"`
	JWT      JWTConfig    `yaml:"jwt"`
	Database DatabaseConfig
}

type ServerConfig struct {
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type JWTConfig struct {
	Secret string        `yaml:"secret"`
	Expiry time.Duration `yaml:"expiry"`
}

type DatabaseConfig struct {
	URL string
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	if port := os.Getenv("PORT"); port != "" {
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil && p > 0 {
			cfg.Server.Port = p
		}
	}

	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.JWT.Secret = secret
	}

	if expiry := os.Getenv("JWT_EXPIRY"); expiry != "" {
		d, err := time.ParseDuration(expiry)
		if err != nil {
			return nil, fmt.Errorf("parse JWT_EXPIRY: %w", err)
		}
		cfg.JWT.Expiry = d
	}

	cfg.Database.URL = os.Getenv("DATABASE_URL")

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 15 * time.Second
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 15 * time.Second
	}
	if cfg.JWT.Expiry == 0 {
		cfg.JWT.Expiry = 24 * time.Hour
	}
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT secret is required (set jwt.secret in config or JWT_SECRET env)")
	}

	return &cfg, nil
}
