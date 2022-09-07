package ampq

import (
	"context"
	"fmt"
	"os"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/pkg/logger"
	ampq "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	config  configs.AMPQ
	conn    *ampq.Connection
	channel *ampq.Channel
	logger  *logger.Logger
}

func New(conf configs.AMPQ, log *logger.Logger) *RabbitClient {
	return &RabbitClient{config: conf, logger: log}
}

func (rb *RabbitClient) Connect() error {
	configAMPQ := rb.config.GetAMPQConfig()
	conn, err := ampq.Dial(configAMPQ.Host + ":" + configAMPQ.Port)
	if err != nil {
		return err
	}

	rb.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	rb.channel = ch

	return nil
}

func (rb *RabbitClient) Close() error {
	if err := rb.channel.Close(); err != nil {
		return err
	}

	if err := rb.conn.Close(); err != nil {
		return err
	}

	return nil
}

func (rb *RabbitClient) Send(ctx context.Context, message string) error {
	queue, err := rb.ConfigureQueue()
	if err != nil {
		return err
	}

	err = rb.channel.PublishWithContext(ctx, "", queue.Name, false, false,
		ampq.Publishing{ContentType: rb.config.GetAMPQConfig().ContentType, Body: []byte(message)})
	if err != nil {
		return err
	}

	rb.logger.Info("[x] Sent successfully message: " + message)
	return nil
}

func (rb *RabbitClient) Receive(ctx context.Context, dst *os.File) error {
	queue, err := rb.ConfigureQueue()
	if err != nil {
		return err
	}

	messages, err := rb.channel.Consume(queue.Name, "",
		true, false, false, false, nil)
	if err != nil {
		return err
	}

	errs := make(chan error, 1)

	go func() {
		for m := range messages {
			// rb.logger.Info("Received: " + string(m.Body))
			_, err = fmt.Fprintf(dst, "Sent notification: %s\n", m.Body)
			if err != nil {
				errs <- err
				return
			}
		}
	}()

	// rb.logger.Info("[*] Waiting for new messages...")

	select {
	case <-ctx.Done():
		break
	case err = <-errs:
		return err
	}

	return nil
}

func (rb *RabbitClient) ConfigureQueue() (ampq.Queue, error) {
	return rb.channel.QueueDeclare(
		rb.config.GetAMPQConfig().QueueName, false, false, false, false, nil)
}
