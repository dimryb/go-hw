package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

var (
	ErrorReceiveEnd = fmt.Errorf("receive end")
	ErrorSendEnd    = fmt.Errorf("send end")
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *telnetClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	dialer := &net.Dialer{}

	conn, err := dialer.DialContext(ctx, "tcp", t.address)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", t.address, err)
	}

	t.conn = conn
	return nil
}

func (t *telnetClient) Send() error {
	if t.conn == nil {
		return fmt.Errorf("connection is not established")
	}
	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		return fmt.Errorf("%w, %w", ErrorSendEnd, err)
	}
	return err
}

func (t *telnetClient) Receive() error {
	if t.conn == nil {
		return fmt.Errorf("connection is not established")
	}
	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		return fmt.Errorf("%w, %w", ErrorReceiveEnd, err)
	}
	return nil
}

func (t *telnetClient) Close() error {
	if t.conn != nil {
		err := t.conn.Close()
		if err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}
	return nil
}
