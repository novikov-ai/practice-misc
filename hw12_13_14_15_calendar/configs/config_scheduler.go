package configs

import (
	"log"
	"os"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs/groups"

	"github.com/BurntSushi/toml"
)

type ConfigScheduler struct {
	Application groups.Application
	Logger      groups.Logger
	Server      groups.Server
	Database    groups.Database
	AMQP        groups.AMQP
	Scheduler   groups.Scheduler
}

func (conf *ConfigScheduler) GetLoggerConfig() groups.Logger {
	return conf.Logger
}

func (conf *ConfigScheduler) GetDatabaseConfig() groups.Database {
	return conf.Database
}

func (conf *ConfigScheduler) GetServerConfig() groups.Server {
	return conf.Server
}

func (conf *ConfigScheduler) GetAMPQConfig() groups.AMQP {
	return conf.AMQP
}

func NewConfigScheduler(path string) ConfigScheduler {
	configRaw, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("reading config %s error: %v", path, err)
	}

	config := ConfigScheduler{}
	err = toml.Unmarshal(configRaw, &config)
	if err != nil {
		log.Fatalf("parsing config %s error: %v", path, err)
	}

	return config
}
