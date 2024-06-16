package initialization

import (
	"errors"
	"time"

	"github.com/rs/zerolog"

	"picodb/internal/config"
	"picodb/internal/database/storage"
	"picodb/internal/database/storage/wal"
	"picodb/internal/tools"
)

const defaultFlushingBatchSize = 100
const defaultFlushingBatchTimeout = time.Millisecond * 10
const defaultMaxSegmentSize = 10 << 20
const defaultWALDataDirectory = "./data/picodb/wal"

func CreateWAL(
	cfg *config.WALConfig,
	logger *zerolog.Logger,
	stream chan<- []wal.LogData,
) (storage.WAL, error) {
	flushingBatchSize := defaultFlushingBatchSize
	flushingBatchTimeout := defaultFlushingBatchTimeout
	maxSegmentSize := defaultMaxSegmentSize
	dataDirectory := defaultWALDataDirectory

	if cfg != nil {
		if cfg.FlushingBatchLength != 0 {
			flushingBatchSize = cfg.FlushingBatchLength
		}

		if cfg.FlushingBatchTimeout != 0 {
			flushingBatchTimeout = cfg.FlushingBatchTimeout
		}

		if cfg.MaxSegmentSize != "" {
			size, err := tools.ParseSize(cfg.MaxSegmentSize)
			if err != nil {
				return nil, errors.New("max segment size is incorrect")
			}

			maxSegmentSize = size
		}

		if cfg.DataDirectory != "" {
			dataDirectory = cfg.DataDirectory
		}

		fsReader := wal.NewFSReader(dataDirectory, logger)
		fsWriter := wal.NewFSWriter(dataDirectory, maxSegmentSize, logger)

		return wal.NewWAL(fsWriter, fsReader, stream, flushingBatchTimeout, flushingBatchSize, logger), nil
	}

	return nil, nil
}
