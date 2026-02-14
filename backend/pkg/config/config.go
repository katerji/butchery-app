package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DB     DBConfig
	JWT    JWTConfig
	Server ServerConfig
}

type DBConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"butchery"`
	Password string `env:"DB_PASSWORD" envDefault:"butchery_secret"`
	Name     string `env:"DB_NAME" envDefault:"butchery_db"`
}

func (c DBConfig) DSN() string {
	u := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:     c.Name,
		RawQuery: "sslmode=disable",
	}
	return u.String()
}

type JWTConfig struct {
	Secret          string        `env:"JWT_SECRET" envDefault:"change-me-to-a-secure-random-string"`
	AccessTokenTTL  time.Duration `env:"JWT_ACCESS_TOKEN_TTL" envDefault:"15m"`
	RefreshTokenTTL time.Duration `env:"JWT_REFRESH_TOKEN_TTL" envDefault:"168h"`
}

type ServerConfig struct {
	Port int `env:"SERVER_PORT" envDefault:"8080"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if cfg.JWT.Secret == "" || cfg.JWT.Secret == "change-me-to-a-secure-random-string" {
		return nil, fmt.Errorf("JWT_SECRET must be set to a secure value")
	}
	return cfg, nil
}
