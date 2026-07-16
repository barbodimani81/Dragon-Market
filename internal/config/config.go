package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2" // Correct v2 provider import
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Server   ServerConfig   `koanf:"server"`
	Database DatabaseConfig `koanf:"database"`
}

type ServerConfig struct {
	Port           int           `koanf:"port"`
	Env            string        `koanf:"env"`
	TimeoutSeconds time.Duration `koanf:"timeout_seconds"`
}

type DatabaseConfig struct {
	Host         string `koanf:"host"`
	Port         int    `koanf:"port"`
	User         string `koanf:"user"`
	Password     string `koanf:"password"`
	DBName       string `koanf:"dbname"`
	SSLMode      string `koanf:"ssl_mode"`
	MaxOpenConns int    `koanf:"max_open_conns"`
	MaxIdleConns int    `koanf:"max_idle_conns"`
}

func Load(configPath string) (*Config, error) {
	k := koanf.New(".")

	// 1. Load default configs from YAML file
	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to load file config: %w", err)
	}

	// 2. Load and parse environment variables using env.Opt configuration
	err := k.Load(env.Provider(".", env.Opt{
		Prefix: "APP_",
		TransformFunc: func(key string, value string) (string, any) {
			// Strip prefix "APP_", lowercase the key, and swap double underscores "__" for "."
			cleanKey := strings.ReplaceAll(
				strings.ToLower(strings.TrimPrefix(key, "APP_")),
				"__",
				".",
			)
			return cleanKey, value
		},
	}), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load environment config: %w", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Convert raw integer config to time.Duration
	cfg.Server.TimeoutSeconds = cfg.Server.TimeoutSeconds * time.Second

	return &cfg, nil
}
