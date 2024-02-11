package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VladPetriv/d2j/config"
	"github.com/VladPetriv/d2j/internal/controller"
	"github.com/VladPetriv/d2j/internal/service"
	"github.com/VladPetriv/d2j/pkg/caching"
	"github.com/VladPetriv/d2j/pkg/database"
	"github.com/VladPetriv/d2j/pkg/encryption"
	"github.com/VladPetriv/d2j/pkg/hashing"
	"github.com/VladPetriv/d2j/pkg/httpserver"
	"github.com/VladPetriv/d2j/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Run runs entire application
func Run(config config.Config, logger logger.Logger) {
	postgreSQL := database.NewPostgreSQL(logger)

	encryptor := encryption.New()
	hasher := hashing.NewBcrypt()

	redis := caching.NewRedis(caching.ConnectionOptions{
		Host:     config.Redis.Host,
		Password: config.Redis.Password,
		Database: config.Redis.Database,
	})

	serviceOptions := service.Options{
		Logger:    logger,
		Config:    config,
		Cacher:    redis,
		Encryptor: encryptor,
		Hasher:    hasher,
		Database:  postgreSQL,
	}

	services := service.Services{
		Database: service.NewDatabaseService(&serviceOptions),
	}

	httpHandler := gin.New()

	controller.New(controller.Options{
		Handler:  httpHandler,
		Logger:   logger,
		Config:   config,
		Services: services,
	})

	httpServer := httpserver.New(
		httpHandler,
		httpserver.Port(config.HTTP.Port),
		httpserver.ReadTimeout(time.Second*60),
		httpserver.WriteTimeout(time.Second*60),
		httpserver.ShutdownTimeout(time.Second*30),
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())

	case err := <-httpServer.Notify():
		logger.Error("app - Run - httpServer.Notify", "err", err)
	}

	err := httpServer.Shutdown()
	if err != nil {
		logger.Error("app - Run - httpServer.Shutdown", "err", err)
	}

	err = redis.Close()
	if err != nil {
		logger.Error("close redis connection", "err", err)
	}
}
