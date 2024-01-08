package compute_test

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	appContext "picodb/internal/context"
	"picodb/internal/database/compute"
)

func TestAnalyzeQuery(t *testing.T) {
	tests := map[string]struct {
		tokens []string
		query  compute.Query
		err    error
	}{
		"empty tokens": {
			tokens: []string{},
			err:    compute.ErrInvalidCommand,
		},
		"invalid command": {
			tokens: []string{"TRUNCATE"},
			err:    compute.ErrInvalidCommand,
		},
		"invalid number arguments for set query": {
			tokens: []string{"SET", "key"},
			err:    compute.ErrInvalidArguments,
		},
		"invalid number arguments for get query": {
			tokens: []string{"GET", "key", "value"},
			err:    compute.ErrInvalidArguments,
		},
		"invalid number arguments for del query": {
			tokens: []string{"GET", "key", "value"},
			err:    compute.ErrInvalidArguments,
		},
		"valid set query": {
			tokens: []string{"SET", "key", "value"},
			query:  compute.NewQuery(compute.SetCommandID, []string{"key", "value"}),
		},
		"valid get query": {
			tokens: []string{"GET", "key"},
			query:  compute.NewQuery(compute.GetCommandID, []string{"key"}),
		},
		"valid del query": {
			tokens: []string{"DEL", "key"},
			query:  compute.NewQuery(compute.DelCommandID, []string{"key"}),
		},
	}

	ctx := appContext.WithTxID(context.Background(), 555)
	logger := zerolog.Nop()
	analyzer := compute.NewAnalyzer(&logger)

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			query, err := analyzer.AnalyzeQuery(ctx, test.tokens)
			require.Equal(t, test.query, query)
			require.Equal(t, test.err, err)
		})
	}
}
