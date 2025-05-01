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
	address string,
	timeout time.Duration,
	clientCloseDelay time.Duration,
	in io.ReadCloser,
	out io.Writer,
) error {
	client := NewTelnetClient(address, timeout, in, out)

	if err := client.Connect(); err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer func() {
		time.Sleep(clientCloseDelay)
		if err := client.Close(); err != nil {
			log.Fatalf("failed to close client: %v", err)
		}
	}()

	sendErrCh := make(chan error)
	receiveErrCh := make(chan error)

	go func() {
		defer close(sendErrCh)
		sendErrCh <- client.Send()
	}()

	go func() {
		defer close(receiveErrCh)
		receiveErrCh <- client.Receive()
	}()

	for {
		select {
		case err := <-sendErrCh:
			return err
		case err := <-receiveErrCh:
			return err
		case <-ctx.Done():
			return nil
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
	err := runTelnetClient(ctx, address, *timeout, 0, os.Stdin, os.Stdout)
	if err != nil {
		if errors.Is(err, ErrorReceiveEnd) {
			log.Print("Connection was canceled from server")
		} else {
			log.Print("Failed: ", err.Error())
		}
	} else {
		log.Print("Connection was canceled.")
	}
	cancel()
}
