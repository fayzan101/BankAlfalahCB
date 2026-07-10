package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	JWT        JWTConfig        `yaml:"jwt"`
	OpenAI     OpenAIConfig     `yaml:"openai"`
	RateLimit  RateLimitConfig  `yaml:"rate_limit"`
	Security   SecurityConfig   `yaml:"security"`
	Database   DatabaseConfig
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

type OpenAIConfig struct {
	APIKey             string        `yaml:"api_key"`
	Model              string        `yaml:"model"`
	MaxTokens          int           `yaml:"max_tokens"`
	Timeout            time.Duration `yaml:"timeout"`
	MaxHistoryMessages int           `yaml:"max_history_messages"`
}

type RateLimitConfig struct {
	IPRequestsPerMinute   int `yaml:"ip_requests_per_minute"`
	UserRequestsPerMinute int `yaml:"user_requests_per_minute"`
}

type SecurityConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
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

	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		cfg.OpenAI.APIKey = apiKey
	}
	if model := os.Getenv("OPENAI_MODEL"); model != "" {
		cfg.OpenAI.Model = model
	}

	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		parts := strings.Split(origins, ",")
		cfg.Security.AllowedOrigins = make([]string, 0, len(parts))
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				cfg.Security.AllowedOrigins = append(cfg.Security.AllowedOrigins, trimmed)
			}
		}
	}

	if v := os.Getenv("RATE_LIMIT_IP_PER_MINUTE"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("parse RATE_LIMIT_IP_PER_MINUTE: %w", err)
		}
		cfg.RateLimit.IPRequestsPerMinute = n
	}
	if v := os.Getenv("RATE_LIMIT_USER_PER_MINUTE"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("parse RATE_LIMIT_USER_PER_MINUTE: %w", err)
		}
		cfg.RateLimit.UserRequestsPerMinute = n
	}

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

	if cfg.OpenAI.Model == "" {
		cfg.OpenAI.Model = "gpt-4o-mini"
	}
	if cfg.OpenAI.MaxTokens == 0 {
		cfg.OpenAI.MaxTokens = 512
	}
	if cfg.OpenAI.Timeout == 0 {
		cfg.OpenAI.Timeout = 30 * time.Second
	}
	if cfg.OpenAI.MaxHistoryMessages == 0 {
		cfg.OpenAI.MaxHistoryMessages = 20
	}

	if cfg.RateLimit.IPRequestsPerMinute == 0 {
		cfg.RateLimit.IPRequestsPerMinute = 60
	}
	if cfg.RateLimit.UserRequestsPerMinute == 0 {
		cfg.RateLimit.UserRequestsPerMinute = 120
	}
	if len(cfg.Security.AllowedOrigins) == 0 {
		cfg.Security.AllowedOrigins = []string{"*"}
	}

	return &cfg, nil
}
