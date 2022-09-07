package configs

import (
	"log"
	"os"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs/groups"

	"github.com/BurntSushi/toml"
)

type ConfigSender struct {
	Application groups.Application
	Logger      groups.Logger
	AMQP        groups.AMQP
}

func (conf *ConfigSender) GetLoggerConfig() groups.Logger {
	return conf.Logger
}

func (conf *ConfigSender) GetAMQPConfig() groups.AMQP {
	return conf.AMQP
}

func NewConfigSender(path string) ConfigSender {
	configRaw, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("reading config %s error: %v", path, err)
	}

	config := ConfigSender{}
	err = toml.Unmarshal(configRaw, &config)
	if err != nil {
		log.Fatalf("parsing config %s error: %v", path, err)
	}

	return config
}
