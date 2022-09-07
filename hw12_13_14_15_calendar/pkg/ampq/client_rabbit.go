package ampq

import (
	"context"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/pkg/logger"
	ampq "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	config  configs.ConfigScheduler
	conn    *ampq.Connection
	channel *ampq.Channel
	logger  *logger.Logger
}

func New(conf configs.ConfigScheduler, log *logger.Logger) *RabbitClient {
	return &RabbitClient{config: conf, logger: log}
}

func (rb *RabbitClient) Connect() error {
	conn, err := ampq.Dial(rb.config.AMQP.Host + ":" + rb.config.AMQP.Port)
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
		ampq.Publishing{ContentType: rb.config.AMQP.ContentType, Body: []byte(message)})
	if err != nil {
		return err
	}

	rb.logger.Info("[x] Sent successfully message: " + message)
	return nil
}

func (rb *RabbitClient) Receive(ctx context.Context) error {
	queue, err := rb.ConfigureQueue()
	if err != nil {
		return err
	}

	messages, err := rb.channel.Consume(queue.Name, "",
		true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for m := range messages {
			rb.logger.Info("Received: " + string(m.Body))
		}
	}()

	rb.logger.Info("[*] Waiting for new messages...")

	<-ctx.Done()
	return nil
}

func (rb *RabbitClient) ConfigureQueue() (ampq.Queue, error) {
	return rb.channel.QueueDeclare(
		rb.config.AMQP.QueueName, false, false, false, false, nil)
}
