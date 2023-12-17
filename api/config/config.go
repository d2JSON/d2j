package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP   HTTP
	Logger Logger
}

type (
	HTTP struct {
		Port string `env:"HTTP_PORT" env-default:"8080"`
	}

	Logger struct {
		LogLevel string `env:"LOGGER_LOG_LEVEL" env-default:"debug"`
	}
)

var (
	config Config
	once   sync.Once
)

func Get() Config {
	once.Do(func() {
		if err := cleanenv.ReadEnv(&config); err != nil {
			log.Fatal("read config: ", err)
		}
	})

	return config
}
