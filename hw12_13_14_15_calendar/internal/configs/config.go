package configs

import (
	"log"
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Application struct {
		Name string `toml:"name"`
	}

	Logger LoggerConfig
	Server ServerConfig

	Database struct {
		InMemory bool   `toml:"in_memory"`
		Driver   string `toml:"driver"`
		Source   string `toml:"source"`
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
