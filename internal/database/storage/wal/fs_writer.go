package wal

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var now = time.Now

type FSWriter struct {
	segment   *os.File
	directory string

	segmentSize    int
	maxSegmentSize int

	logger *zerolog.Logger
}

func NewFSWriter(directory string, maxSegmentSize int, logger *zerolog.Logger) *FSWriter {
	return &FSWriter{
		directory:      directory,
		maxSegmentSize: maxSegmentSize,
		logger:         logger,
	}
}

func (w *FSWriter) WriteBatch(batch []Log) {
	if w.segment == nil {
		if err := w.rotateSegment(); err != nil {
			w.acknowledgeWrite(batch, err)

			return
		}
	}

	if w.segmentSize > w.maxSegmentSize {
		if err := w.rotateSegment(); err != nil {
			w.acknowledgeWrite(batch, err)

			return
		}
	}

	logs := make([]LogData, 0, len(batch))
	for _, log := range batch {
		logs = append(logs, log.data)
	}

	if err := w.writeLogs(logs); err != nil {
		w.acknowledgeWrite(batch, err)

		return
	}

	err := w.segment.Sync()
	if err != nil {
		w.logger.Error().Err(err).Msg("failed to sync segment file")
	}

	w.acknowledgeWrite(batch, err)
}

func (w *FSWriter) writeLogs(logs []LogData) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(logs); err != nil {
		w.logger.Warn().Err(err).Msg("failed to encode logs data")

		return err
	}

	writtenBytes, err := w.segment.Write(buffer.Bytes())
	if err != nil {
		w.logger.Warn().Err(err).Msg("failed to write logs data")

		return err
	}

	w.segmentSize += writtenBytes

	return nil
}

func (w *FSWriter) acknowledgeWrite(batch []Log, err error) {
	for _, log := range batch {
		log.SetResult(err)
	}
}

func (w *FSWriter) rotateSegment() error {
	segmentName := fmt.Sprintf("%s/wal_%d.log", w.directory, now().UnixMilli())

	flags := os.O_CREATE | os.O_WRONLY
	segment, err := os.OpenFile(segmentName, flags, 0o644)
	if err != nil {
		w.logger.Error().Err(err).Msg("failed to create wal segment")

		return err
	}

	w.segment = segment
	w.segmentSize = 0

	return nil
}
