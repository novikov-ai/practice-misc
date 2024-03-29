package ampq

import (
	"context"
	"os"
)

type Client interface {
	Connect() error
	Close() error
	Send(ctx context.Context, message string) error
	Receive(ctx context.Context, dst *os.File) error
}
