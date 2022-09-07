package configs

import "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs/groups"

type Configurator interface {
	Logger
	GetDatabaseConfig() groups.Database
	GetServerConfig() groups.Server
}

type Logger interface {
	GetLoggerConfig() groups.Logger
}

type AMPQ interface {
	GetLoggerConfig() groups.Logger
	GetAMPQConfig() groups.AMQP
}
