package compute_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"picodb/internal/database/compute"

	appContext "picodb/internal/context"
)

func TestParse(t *testing.T) {
	tests := map[string]struct {
		query  string
		tokens []string
		err    error
	}{
		"empty query": {
			query: "",
		},
		"query without tokens": {
			query: "   ",
		},
		"query with UTF symbols": {
			query: "字文下",
			err:   compute.ErrInvalidSymbol,
		},
		"query with one token": {
			query:  "set",
			tokens: []string{"set"},
		},
		"query with two tokens": {
			query:  "set key",
			tokens: []string{"set", "key"},
		},
		"query with one token with digits": {
			query:  "2set1",
			tokens: []string{"2set1"},
		},
		"query with one token with underscores": {
			query:  "_set__",
			tokens: []string{"_set__"},
		},
		"query with one token with invalid symbols": {
			query: ".set#",
			err:   compute.ErrInvalidSymbol,
		},
		"query with two tokens with additional spaces": {
			query:  " set   key  ",
			tokens: []string{"set", "key"},
		},
	}

	ctx := appContext.WithTxID(context.Background(), 555)

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			logger := zerolog.Nop()
			parser := compute.NewParser(&logger)

			tokens, err := parser.ParseQuery(ctx, test.query)
			require.Equal(t, test.err, err)
			require.True(t, reflect.DeepEqual(test.tokens, tokens))
		})
	}
}
