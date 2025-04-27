package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}

func TestRunTelnetClient(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, listener.Close())
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			return
		default:
			conn, err := listener.Accept()
			defer func() {
				require.NoError(t, conn.Close())
			}()
			if err != nil {
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					return
				}
				t.Errorf("Failed to accept connection: %v", err)
				return
			}

			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil && ctx.Err() == nil {
				t.Errorf("Failed to read from client: %v", err)
				return
			}
			clientData := string(buf[:n])
			require.Equal(t, "Hello\nFrom\nNC\n", clientData)

			_, err = conn.Write([]byte("I\nam\nTELNET client\n"))
			if err != nil && ctx.Err() == nil {
				t.Errorf("Failed to write to client: %v", err)
				return
			}
		}
	}()

	serverAddr := listener.Addr().String()

	input := strings.NewReader("Hello\nFrom\nNC\n")
	output := &bytes.Buffer{}

	err = runTelnetClient(serverAddr, 5*time.Second, io.NopCloser(input), output)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Errorf("Client timed out")
	}
	require.NoError(t, err)

	expectedOutput := "I\nam\nTELNET client\n"
	require.Equal(t, expectedOutput, output.String())

	time.Sleep(100 * time.Millisecond)

	wg.Wait()
}
