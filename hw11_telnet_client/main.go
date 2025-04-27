package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func runTelnetClient(
	ctx context.Context,
	address string, timeout time.Duration,
	in io.ReadCloser, out io.Writer,
) error {
	client := NewTelnetClient(address, timeout, in, out)

	if err := client.Connect(); err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalf("failed to close client: %v", err)
		}
	}()

	sendErrCh := make(chan error)
	go func() {
		sendErrCh <- client.Send()
		close(sendErrCh)
	}()

	receiveErrCh := make(chan error)
	go func() {
		receiveErrCh <- client.Receive()
		close(receiveErrCh)
	}()

	for {
		select {
		case err := <-sendErrCh:
			if err != nil && !errors.Is(err, io.EOF) {
				return fmt.Errorf("send error: %w", err)
			}
		case err := <-receiveErrCh:
			if err != nil && !errors.Is(err, io.EOF) {
				return fmt.Errorf("receive error: %w", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "Timeout for connection")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Println("Usage: go-telnet [--timeout=TIMEOUT] host port")
		os.Exit(1)
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	ctx, cancel := context.WithCancel(context.Background())
	err := runTelnetClient(ctx, address, *timeout, os.Stdin, os.Stdout)
	log.Printf("Error: %v\n", err)

	cancel()
	os.Exit(1)
}
