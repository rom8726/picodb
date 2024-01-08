package compute_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"picodb/internal/database/compute"
)

func TestCommandNameToCommandID(t *testing.T) {
	require.Equal(t, compute.SetCommandID, compute.CommandNameToCommandID("SET"))
	require.Equal(t, compute.GetCommandID, compute.CommandNameToCommandID("GET"))
	require.Equal(t, compute.DelCommandID, compute.CommandNameToCommandID("DEL"))
	require.Equal(t, compute.UnknownCommandID, compute.CommandNameToCommandID("TRUNCATE"))
}
