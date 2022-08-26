package main

import (
	"errors"
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

type Client struct {
	timeout time.Duration
	input   io.ReadCloser
	output  io.Writer
	conn    net.Conn
	address string
}

func (cl *Client) Connect() (err error) {
	cl.conn, err = net.DialTimeout("tcp", cl.address, cl.timeout)
	return err
}

func (cl *Client) Close() error {
	if cl.conn != nil {
		return cl.conn.Close()
	}

	return errors.New("connection is closed already")
}

func (cl *Client) Send() error {
	_, err := io.Copy(cl.conn, cl.input)
	return err
}

func (cl *Client) Receive() error {
	_, err := io.Copy(cl.output, cl.conn)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{address: address, timeout: timeout, input: in, output: out}
}
