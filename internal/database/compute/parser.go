package compute

import (
	"context"

	"github.com/rs/zerolog"

	appContext "picodb/internal/context"
)

type Parser struct {
	logger *zerolog.Logger
}

func NewParser(logger *zerolog.Logger) *Parser {
	return &Parser{
		logger: logger,
	}
}

func (p *Parser) ParseQuery(ctx context.Context, query string) ([]string, error) {
	machine := newStateMachine()
	tokens, err := machine.parse(query)
	if err != nil {
		return nil, err
	}

	if p.logger.GetLevel() == zerolog.DebugLevel {
		p.logger.Debug().
			Int64("tx", appContext.TxIDFromContext(ctx)).
			Strs("tokens", tokens).
			Msg("query parsed")
	}

	return tokens, nil
}

func isWhiteSpace(symbol byte) bool {
	return symbol == '\t' || symbol == '\n' || symbol == ' '
}

func isLetter(symbol byte) bool {
	return (symbol >= 'a' && symbol <= 'z') ||
		(symbol >= 'A' && symbol <= 'Z') ||
		(symbol >= '0' && symbol <= '9') ||
		(symbol == '_')
}
