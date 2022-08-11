package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
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
		stderrLog(ErrTimeout.Error())
		return
	}
	defer tc.Close()

	stderrLog(fmt.Sprintf("Connected to %s\n", hostPort))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		clientDo(ctx, tc.Send)
		stop()
	}()
	go func() {
		clientDo(ctx, tc.Receive)
	}()

	<-ctx.Done()
}

func stderrLog(s string) {
	_, err := fmt.Fprintf(os.Stderr, "%s\n", s)
	if err != nil {
		log.Fatal(err)
	}
}

func getHostAddress() (string, error) {
	arguments := os.Args
	for i, arg := range arguments {
		if strings.Contains(arg, FlagTimeout) {
			arguments = append(arguments[0:i], arguments[i+1:]...)
			break
		}
	}

	if len(arguments) < 3 {
		return "", ErrWrongAddress
	}

	return fmt.Sprintf("%s:%s", arguments[1], arguments[2]), nil
}

func clientDo(ctx context.Context, work func() error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := work()
			if errors.Is(io.EOF, err) {
				stderrLog(err.Error())
				return
			} else if err != nil {
				log.Fatal(err)
			}
		}
	}
}
