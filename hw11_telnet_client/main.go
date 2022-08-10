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
	Timeout    time.Duration
	ErrTimeout = errors.New("timeout reached")
)

func init() {
	flag.DurationVar(&Timeout, FlagTimeout, 10*time.Second, "timeout for connection")
}

func main() {
	flag.Parse()

	hostPort := getHostAddress()

	tc := NewTelnetClient(hostPort, Timeout, os.Stdin, os.Stdout)
	defer tc.Close()

	err := tc.Connect()
	if err != nil {
		log.Fatal(ErrTimeout)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		telnetClientDo(ctx, tc.Send)
		stop()
	}()
	go func() {
		telnetClientDo(ctx, tc.Receive)
		stop()
	}()

	select {
	case <-ctx.Done():
		fmt.Println()
		stop()
	}
}

func logError(r error) {
	_, err := fmt.Fprintf(os.Stderr, "%s\n", r)
	if err != nil {
		log.Fatal(err)
	}
}

func getHostAddress() string {
	arguments := os.Args
	for i, arg := range arguments {
		if strings.Contains(arg, FlagTimeout) {
			left := arguments[0:i]
			right := arguments[i+1:]
			arguments = append(left, right...)
			break
		}
	}

	if len(arguments) < 3 {
		log.Fatal("Please provide host and port")
	}

	return fmt.Sprintf("%s:%s", arguments[1], arguments[2])
}

func telnetClientDo(ctx context.Context, work func() error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := work()
			if errors.Is(io.EOF, err) {
				logError(err)
				return
			} else if err != nil {
				log.Fatal(err)
			}
		}
	}
}
