package inmemory

import (
	"context"
	"errors"
	"hash/fnv"

	"github.com/rs/zerolog"

	appContext "picodb/internal/context"
	"picodb/internal/database/compute"
	"picodb/internal/database/storage/wal"
)

type hashTable interface {
	Set(string, string)
	Get(string) (string, bool)
	Del(string)
}

type Engine struct {
	partitions []hashTable
	logger     *zerolog.Logger
}

func NewEngine(
	tableBuilder func() hashTable,
	stream <-chan []wal.LogData,
	partitionsNumber int,
	logger *zerolog.Logger,
) (*Engine, error) {
	if tableBuilder == nil {
		return nil, errors.New("hash table builder is invalid")
	}

	if partitionsNumber <= 0 {
		return nil, errors.New("partitions number is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	partitions := make([]hashTable, partitionsNumber)
	for i := 0; i < partitionsNumber; i++ {
		if partition := tableBuilder(); partition != nil {
			partitions[i] = partition
		} else {
			return nil, errors.New("hash table partition is invalid")
		}
	}

	engine := &Engine{
		partitions: partitions,
		logger:     logger,
	}

	if stream != nil {
		go func() {
			for logs := range stream {
				engine.applyLogs(logs)
			}
		}()
	}

	return engine, nil
}

func (e *Engine) Set(ctx context.Context, key, value string) {
	idx := e.partitionIdx(key)
	partition := e.partitions[idx]
	partition.Set(key, value)

	txID := appContext.TxIDFromContext(ctx)
	e.logger.Debug().Int64("tx", txID).Msg("success set query")
}

func (e *Engine) Get(ctx context.Context, key string) (string, bool) {
	idx := e.partitionIdx(key)
	partition := e.partitions[idx]
	value, found := partition.Get(key)

	txID := appContext.TxIDFromContext(ctx)
	e.logger.Debug().Int64("tx", txID).Msg("success get query")

	return value, found
}

func (e *Engine) Del(ctx context.Context, key string) {
	idx := e.partitionIdx(key)
	partition := e.partitions[idx]
	partition.Del(key)

	txID := appContext.TxIDFromContext(ctx)
	e.logger.Debug().Int64("tx", txID).Msg("success del query")
}

func (e *Engine) partitionIdx(key string) int {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(key))

	return int(hash.Sum32()) % len(e.partitions)
}

func (e *Engine) applyLogs(logs []wal.LogData) {
	for _, log := range logs {
		switch log.CommandID {
		case compute.SetCommandID:
			e.Set(context.Background(), log.Arguments[0], log.Arguments[1])
		case compute.DelCommandID:
			e.Del(context.Background(), log.Arguments[0])
		}
	}
}
