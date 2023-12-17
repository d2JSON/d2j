package main

import (
	"github.com/VladPetriv/postgreSQL2JSON/config"
	"github.com/VladPetriv/postgreSQL2JSON/internal/app"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/logger"
)

func main() {
	config := config.Get()
	logger := logger.New(config.Logger.LogLevel)

	app.Run(config, logger)
}
