//nolint:gocritic
package initialization

import (
	"errors"
	"os"

	"github.com/rs/zerolog"

	"picodb/internal/config"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

var supportedLoggingLevels = map[string]zerolog.Level{
	DebugLevel: zerolog.DebugLevel,
	InfoLevel:  zerolog.InfoLevel,
	WarnLevel:  zerolog.WarnLevel,
	ErrorLevel: zerolog.ErrorLevel,
}

// const defaultEncoding = "json"
const defaultLevel = zerolog.InfoLevel

// const defaultOutputPath = "output.log"

func CreateLogger(cfg *config.LoggingConfig) (*zerolog.Logger, error) {
	level := defaultLevel
	// output := defaultOutputPath

	if cfg != nil {
		if cfg.Level != "" {
			var found bool
			if level, found = supportedLoggingLevels[cfg.Level]; !found {
				return nil, errors.New("logging level is incorrect")
			}
		}

		// if cfg.Output != "" {
		//	output = cfg.Output
		// }
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger().Level(level)

	return &logger, nil
}
