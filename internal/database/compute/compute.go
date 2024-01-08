//go:generate mockery --name QueryParser --case snake
//go:generate mockery --name QueryAnalyzer --case snake
package compute

import (
	"context"

	"github.com/rs/zerolog"
)

type QueryParser interface {
	ParseQuery(context.Context, string) ([]string, error)
}

type QueryAnalyzer interface {
	AnalyzeQuery(context.Context, []string) (Query, error)
}

type Compute struct {
	parser   QueryParser
	analyzer QueryAnalyzer
	logger   *zerolog.Logger
}

func NewCompute(parser QueryParser, analyzer QueryAnalyzer, logger *zerolog.Logger) *Compute {
	return &Compute{
		parser:   parser,
		analyzer: analyzer,
		logger:   logger,
	}
}

func (d *Compute) HandleQuery(ctx context.Context, queryStr string) (Query, error) {
	tokens, err := d.parser.ParseQuery(ctx, queryStr)
	if err != nil {
		return Query{}, err
	}

	query, err := d.analyzer.AnalyzeQuery(ctx, tokens)
	if err != nil {
		return Query{}, err
	}

	return query, nil
}
