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
	tests := []struct {
		name           string
		inputData      string
		expectedOutput string
	}{
		{
			name:           "With input data",
			inputData:      "Hello\nFrom\nNC\n",
			expectedOutput: "I\nam\nTELNET client\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
				require.Equal(t, tt.inputData, clientData)

				_, err = conn.Write([]byte("I\nam\nTELNET client\n"))
				require.NoError(t, err)
			}()

			serverAddr := listener.Addr().String()

			input := &SafeBuffer{}
			output := &SafeBuffer{}

			ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
			defer cancel()

			input.Write([]byte(tt.inputData))

			wg.Add(1)
			go func() {
				defer wg.Done()
				err := runTelnetClient(ctx, serverAddr, 5*time.Second, 100*time.Millisecond, io.NopCloser(input), output)
				if err != nil && !errors.Is(err, context.DeadlineExceeded) {
					t.Errorf("Client error: %v", err)
				}
			}()

			wg.Wait()

			select {
			case <-ctx.Done():
				t.Fatal("Test timed out")
			default:
				require.Equal(t, tt.expectedOutput, output.String())
			}
		})
	}
}

type SafeBuffer struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func (sb *SafeBuffer) Write(p []byte) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.Write(p)
}

func (sb *SafeBuffer) Read(p []byte) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.Read(p)
}

func (sb *SafeBuffer) String() string {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.String()
}
