package main

import (
	"bufio"
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
	return cl.conn.Close()
}

func (cl *Client) Send() error {
	dataSent := false
	r := bufio.NewReader(cl.input)
	for {
		line, err := r.ReadString('\n')
		if errors.Is(err, io.EOF) {
			if !dataSent {
				return io.EOF
			}

			break
		}
		_, err = cl.conn.Write([]byte(line))
		if err != nil {
			return err
		}
		dataSent = true
	}
	return nil
}

func (cl *Client) Receive() error {
	_, err := io.Copy(cl.output, cl.conn)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{address: address, timeout: timeout, input: in, output: out}
}
