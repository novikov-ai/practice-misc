package configs

import (
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/logger"
	"github.com/pelletier/go-toml"
	"log"
	"os"
)

type Config struct {
	Logger logger.Config

	Application struct {
		Name string `toml:"name"`
	}

	Server struct {
		ServerHost string `toml:"host"`
		ServerPort string `toml:"port"`
	}

	Database struct {
		DbModeInMemory bool   `toml:"in_memory"`
		DbConnection   string `toml:"connection"`
	}
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
