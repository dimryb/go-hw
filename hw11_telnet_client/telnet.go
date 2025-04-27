package main

import (
	"context"
	"io"
	"net"
	"time"
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

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.

func (t *telnetClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	dialer := &net.Dialer{}

	conn, err := dialer.DialContext(ctx, "tcp", t.address)
	if err != nil {
		return err
	}

	t.conn = conn
	return nil
}

func (t *telnetClient) Send() error {
	return nil
}

func (t *telnetClient) Receive() error {
	return nil
}

func (t *telnetClient) Close() error {
	return t.in.Close()
}
