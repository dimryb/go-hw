package main

import (
	"io"
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
