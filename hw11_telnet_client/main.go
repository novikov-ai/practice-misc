package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	FlagTimeout = "timeout"
)

var (
	Timeout         time.Duration
	ErrTimeout      = errors.New("timeout reached")
	ErrWrongAddress = errors.New("incorrect host or port")
)

func init() {
	flag.DurationVar(&Timeout, FlagTimeout, 10*time.Second, "timeout for connection")
}

func main() {
	flag.Parse()

	hostPort, err := getHostAddress()
	if err != nil {
		fmt.Println("Please provide host and port")
		return
	}

	tc := NewTelnetClient(hostPort, Timeout, os.Stdin, os.Stdout)
	if tc.Connect() != nil {
		fmt.Fprintln(os.Stderr, ErrTimeout.Error())
		return
	}
	defer tc.Close()

	fmt.Fprintln(os.Stderr, "Connected to", hostPort)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		err = tc.Send()
		if err != nil {
			fmt.Println("send error: ", err)
		}

		fmt.Fprintln(os.Stderr, "...EOF")
		stop()
	}()
	go func() {
		err = tc.Receive()
		if err != nil {
			fmt.Println("receive error:", err)
		}
		fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
		stop()
	}()

	<-ctx.Done()
}

func getHostAddress() (string, error) {
	arguments := flag.Args()

	if len(arguments) < 2 {
		return "", ErrWrongAddress
	}

	return fmt.Sprintf("%s:%s", arguments[0], arguments[1]), nil
}
