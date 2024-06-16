package initialization

import (
	"errors"

	"github.com/rs/zerolog"

	"picodb/internal/config"
	"picodb/internal/database/storage"
	inMemory "picodb/internal/database/storage/engine/in-memory"
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
		return inMemory.NewEngine(inMemory.HashTableBuilder, nil, defaultPartitionsNumber, logger)
	}

	if cfg.Type != "" {
		_, found := supportedEngineTypes[cfg.Type]
		if !found {
			return nil, errors.New("engine type is incorrect")
		}
	}

	return inMemory.NewEngine(inMemory.HashTableBuilder, stream, defaultPartitionsNumber, logger)
}
