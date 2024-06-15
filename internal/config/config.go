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
	Engine      *EngineConfig      `yaml:"engine"`
	WAL         *WALConfig         `yaml:"wal"`
	Replication *ReplicationConfig `yaml:"replication"`
	Network     *NetworkConfig     `yaml:"network"`
	Logging     *LoggingConfig     `yaml:"logging"`
}

//nolint:tagliatelle // it's ok
type NetworkConfig struct {
	Address        string        `yaml:"address"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

type EngineConfig struct {
	Type string `yaml:"type"`
}

type WALConfig struct {
	FlushingBatchLength  int           `yaml:"flushing_batch_length"`
	FlushingBatchTimeout time.Duration `yaml:"flushing_batch_timeout"`
	MaxSegmentSize       string        `yaml:"max_segment_size"`
	DataDirectory        string        `yaml:"data_directory"`
}

type ReplicationConfig struct {
	ReplicaType   string        `yaml:"replica_type"`
	MasterAddress string        `yaml:"master_address"`
	SyncInterval  time.Duration `yaml:"sync_interval"`
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
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return Config{}, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}

func validate(cfg *Config) error {
	err := validation.ValidateStruct(cfg.Network,
		validation.Field(&cfg.Network.Address, validation.Required),
	)
	if err != nil {
		return fmt.Errorf("validate network section: %w", err)
	}

	return nil
}
