package groups

type AMQP struct {
	Host        string `toml:"host"`
	Port        string `toml:"port"`
	QueueName   string `toml:"queue_name"`
	ContentType string `toml:"content_type"`
}
