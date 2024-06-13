// Package config is a configuration package
package config

import (
	"fmt"
	"os"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gopkg.in/yaml.v3"
)

const (
	configDefaultFilePath = "config.yml"
)

type Config struct {
	Network NetworkConfig `yaml:"network"`
	Logging LoggingConfig `yaml:"logging"`
}

//nolint:tagliatelle // it's ok
type NetworkConfig struct {
	Port           uint16        `yaml:"port"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

func Init() (Config, error) {
	var configPath string

	if len(os.Args) > 1 {
		configPath = os.Args[1]
	} else {
		configPath = configDefaultFilePath
	}

	info, err := os.Stat(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("stat config %q: %w", configPath, err)
	}

	if info.IsDir() {
		return Config{}, fmt.Errorf("config %q is a directory", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	cfg := Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return Config{}, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}

func validate(cfg *Config) error {
	netw := cfg.Network
	err := validation.ValidateStruct(&netw,
		validation.Field(&netw.Port, validation.Required, validation.Min(uint16(1024))),
	)
	if err != nil {
		return fmt.Errorf("validate network section: %w", err)
	}

	return nil
}
