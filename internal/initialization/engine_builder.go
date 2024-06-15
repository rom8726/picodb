package initialization

import (
	"errors"

	"github.com/rs/zerolog"

	"picodb/internal/config"
	"picodb/internal/database/storage"
	"picodb/internal/database/storage/engine/in-memory"
	"picodb/internal/database/storage/wal"
)

const (
	InMemoryEngine = "in_memory"
)

var supportedEngineTypes = map[string]struct{}{
	InMemoryEngine: {},
}

const defaultPartitionsNumber = 10

func CreateEngine(
	cfg *config.EngineConfig,
	logger *zerolog.Logger,
	stream <-chan []wal.LogData,
) (storage.Engine, error) {
	if cfg == nil {
		return inmemory.NewEngine(inmemory.HashTableBuilder, nil, defaultPartitionsNumber, logger)
	}

	if cfg.Type != "" {
		_, found := supportedEngineTypes[cfg.Type]
		if !found {
			return nil, errors.New("engine type is incorrect")
		}
	}

	return inmemory.NewEngine(inmemory.HashTableBuilder, stream, defaultPartitionsNumber, logger)
}
