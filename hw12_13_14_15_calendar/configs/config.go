package configs

import "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs/groups"

type Configurator interface {
	GetLoggerConfig() groups.Logger
	GetDatabaseConfig() groups.Database
	GetServerConfig() groups.Server
}
