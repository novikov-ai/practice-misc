package sender

import (
	"context"
	"os"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/app"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/pkg/ampq"
)

type Sender struct {
	logger        app.Logger
	config        configs.ConfigSender
	messageBroker ampq.Client
}

func New(log app.Logger, conf configs.ConfigSender, broker ampq.Client) *Sender {
	return &Sender{logger: log, config: conf, messageBroker: broker}
}

func (s *Sender) Start(ctx context.Context, receiver *os.File) error {
	return s.messageBroker.Receive(ctx, receiver)
}
