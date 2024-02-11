package main

import (
	"github.com/VladPetriv/d2j/config"
	"github.com/VladPetriv/d2j/internal/app"
	"github.com/VladPetriv/d2j/pkg/logger"
	_ "github.com/lib/pq"
)

func main() {
	config := config.Get()
	logger := logger.NewSlog(config.Logger.LogLevel)

	app.Run(config, logger)
}
