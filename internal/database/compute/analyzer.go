package compute

import (
	"context"
	"errors"

	"github.com/rs/zerolog"

	appContext "picodb/internal/context"
)

const (
	setQueryArgumentsNumber = 2
	getQueryArgumentsNumber = 1
	delQueryArgumentsNumber = 1
)

var queryArgumentsNumber = map[CommandID]int{
	SetCommandID: setQueryArgumentsNumber,
	GetCommandID: getQueryArgumentsNumber,
	DelCommandID: delQueryArgumentsNumber,
}

var (
	ErrInvalidSymbol    = errors.New("invalid symbol")
	ErrInvalidCommand   = errors.New("invalid command")
	ErrInvalidArguments = errors.New("invalid arguments")
)

type Analyzer struct {
	logger *zerolog.Logger
}

func NewAnalyzer(logger *zerolog.Logger) *Analyzer {
	return &Analyzer{
		logger: logger,
	}
}

func (a *Analyzer) AnalyzeQuery(ctx context.Context, tokens []string) (Query, error) {
	if len(tokens) == 0 {
		return Query{}, ErrInvalidCommand
	}

	command := tokens[0]
	commandID := CommandNameToCommandID(command)
	if commandID == UnknownCommandID {
		return Query{}, ErrInvalidCommand
	}

	query := NewQuery(commandID, tokens[1:])
	argumentsNumber := queryArgumentsNumber[commandID]
	if len(query.Arguments()) != argumentsNumber {
		return Query{}, ErrInvalidArguments
	}

	if a.logger.GetLevel() == zerolog.DebugLevel {
		a.logger.Debug().Int64("tx", appContext.TxIDFromContext(ctx)).Msg("query analyzed")
	}

	return query, nil
}
