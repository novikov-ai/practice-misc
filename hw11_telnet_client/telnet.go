package main

import (
	"bufio"
	"fmt"
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
	address    string
	timeout    time.Duration
	input      io.ReadCloser
	output     io.Writer
	connection net.Conn
}

func (cl *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", cl.address, cl.timeout)
	if err != nil {
		return err
	}

	cl.connection = conn

	fmt.Printf("Welcome to %s!\n", cl.address)
	return nil
}

func (cl *Client) Close() error {
	return cl.connection.Close()
}

func (cl *Client) Send() error {
	r := bufio.NewReader(cl.input)
	message, err := r.ReadString('\n')
	if err != nil {
		return err
	}

	_, err = cl.connection.Write([]byte(message))
	return err
}

func (cl *Client) Receive() error {
	r := bufio.NewReader(cl.connection)

	message, err := r.ReadString('\n')
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(cl.output, message)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{address: address, timeout: timeout, input: in, output: out}
}
