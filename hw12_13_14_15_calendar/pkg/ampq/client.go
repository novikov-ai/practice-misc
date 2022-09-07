package ampq

import "context"

type Client interface {
	Connect() error
	Close() error
	Send(ctx context.Context, message string) error
	Receive(ctx context.Context) error
}
