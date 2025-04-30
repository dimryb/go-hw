package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
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

	go func() {
		defer wg.Done()

		conn, err := listener.Accept()
		defer func() {
			require.NoError(t, conn.Close())
		}()
		require.NoError(t, err)

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		require.NoError(t, err)
		clientData := string(buf[:n])
		require.Equal(t, "Hello\nFrom\nNC\n", clientData)

		_, err = conn.Write([]byte("I\nam\nTELNET client\n"))
		require.NoError(t, err)
	}()

	serverAddr := listener.Addr().String()

	input := &bytes.Buffer{}
	output := &bytes.Buffer{}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		err := runTelnetClient(ctx, serverAddr, 5*time.Second, io.NopCloser(input), output)
		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Client error: %v", err)
		}
	}()

	input.WriteString("Hello\nFrom\nNC\n")

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("Test timed out")
	}

	expectedOutput := "I\nam\nTELNET client\n"
	require.Equal(t, expectedOutput, output.String())

	wg.Wait()
}
