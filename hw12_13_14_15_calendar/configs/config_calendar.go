package configs

import (
	"log"
	"os"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs/groups"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Application groups.Application
	Logger      groups.Logger
	Server      groups.Server
	Database    groups.Database
}

func (conf *Config) GetLoggerConfig() groups.Logger {
	return conf.Logger
}

func (conf *Config) GetDatabaseConfig() groups.Database {
	return conf.Database
}

func (conf *Config) GetServerConfig() groups.Server {
	return conf.Server
}

func NewConfig(path string) Config {
	configRaw, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("reading config %s error: %v", path, err)
	}

	config := Config{}
	err = toml.Unmarshal(configRaw, &config)
	if err != nil {
		log.Fatalf("parsing config %s error: %v", path, err)
	}

	return config
}
