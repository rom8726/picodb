package compute_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"picodb/internal/database/compute"
)

func TestQuery(t *testing.T) {
	query := compute.NewQuery(compute.GetCommandID, []string{"GET", "key"})
	require.Equal(t, compute.GetCommandID, query.CommandID())
	require.True(t, reflect.DeepEqual([]string{"GET", "key"}, query.Arguments()))
}
