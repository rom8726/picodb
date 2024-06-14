package network

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestTCPServer(t *testing.T) {
	t.Parallel()

	request := "hello server"
	response := "hello client"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maxMessageSize := 2048
	maxConnectionsNumber := 10
	idleTimeout := time.Minute
	logger := zerolog.Nop()
	server, err := NewTCPServer(":20001", maxConnectionsNumber, maxMessageSize, idleTimeout, &logger)
	require.NoError(t, err)

	go func() {
		require.NoError(t, server.HandleQueries(ctx, func(ctx context.Context, buffer []byte) []byte {
			require.True(t, reflect.DeepEqual([]byte(request), buffer))
			return []byte(response)
		}))
	}()

	connection, err := net.Dial("tcp", "localhost:20001")
	require.NoError(t, err)

	_, err = connection.Write([]byte(request))
	require.NoError(t, err)

	buffer := make([]byte, 2048)
	count, err := connection.Read(buffer)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual([]byte(response), buffer[:count]))
}
