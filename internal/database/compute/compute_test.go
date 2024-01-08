package compute_test

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	appContext "picodb/internal/context"
	"picodb/internal/database/compute"
	"picodb/internal/database/compute/mocks"
)

func TestHandleQueryWithParsingError(t *testing.T) {
	ctx := appContext.WithTxID(context.Background(), 555)

	parser := mocks.NewQueryParser(t)
	parser.On("ParseQuery", mock.Anything, "## key").Return(nil, compute.ErrInvalidCommand)
	analyzer := mocks.NewQueryAnalyzer(t)

	log := zerolog.Nop()
	comp := compute.NewCompute(parser, analyzer, &log)

	query, err := comp.HandleQuery(ctx, "## key")
	require.Error(t, err, compute.ErrInvalidCommand)
	require.Equal(t, compute.Query{}, query)
}

func TestHandleQueryWithAnalyzingError(t *testing.T) {
	ctx := appContext.WithTxID(context.Background(), 555)

	parser := mocks.NewQueryParser(t)
	parser.On("ParseQuery", mock.Anything, "TRUNCATE key").
		Return([]string{"TRUNCATE", "key"}, nil)
	analyzer := mocks.NewQueryAnalyzer(t)
	analyzer.On("AnalyzeQuery", mock.Anything, []string{"TRUNCATE", "key"}).
		Return(compute.Query{}, compute.ErrInvalidCommand)

	log := zerolog.Nop()
	comp := compute.NewCompute(parser, analyzer, &log)

	query, err := comp.HandleQuery(ctx, "TRUNCATE key")
	require.Error(t, err, compute.ErrInvalidCommand)
	require.Equal(t, compute.Query{}, query)
}

func TestHandleQuery(t *testing.T) {
	ctx := appContext.WithTxID(context.Background(), 555)

	parser := mocks.NewQueryParser(t)
	parser.On("ParseQuery", mock.Anything, "GET key").
		Return([]string{"GET", "key"}, nil)
	analyzer := mocks.NewQueryAnalyzer(t)
	analyzer.On("AnalyzeQuery", mock.Anything, []string{"GET", "key"}).
		Return(compute.NewQuery(compute.GetCommandID, []string{"key"}), nil)

	log := zerolog.Nop()
	comp := compute.NewCompute(parser, analyzer, &log)

	query, err := comp.HandleQuery(ctx, "GET key")
	require.NoError(t, err)
	require.Equal(t, compute.NewQuery(compute.GetCommandID, []string{"key"}), query)
}
