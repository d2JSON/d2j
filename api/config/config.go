package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config represents a structure that contains a configurations for different part of application.
type Config struct {
	HTTP   HTTP
	Logger Logger
	Redis  Redis
}

type (
	// HTTP represents a configuration for HTTP server.
	HTTP struct {
		Port string `env:"PORT" env-default:"8080"`
	}

	// Logger represents a configuration for logger.
	Logger struct {
		LogLevel string `env:"LOGGER_LOG_LEVEL" env-default:"debug"`
	}

	// Redis represents a configuration for Redis client.
	Redis struct {
		Host     string `env:"REDIS_HOST" env-default:"localhost:6379"`
		Password string `env:"REDIS_PASSWORD" env-default:""`
		Database int    `env:"REDIS_DATABASE" env-default:"0"`
	}
)

var (
	config Config
	once   sync.Once
)

// Get returns the config.
func Get() Config {
	once.Do(func() {
		if err := cleanenv.ReadEnv(&config); err != nil {
			log.Fatal("read config: ", err)
		}
	})

	return config
}
