package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/VladPetriv/postgreSQL2JSON/config"
	"github.com/VladPetriv/postgreSQL2JSON/internal/controller"
	"github.com/VladPetriv/postgreSQL2JSON/internal/service"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/caching"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/database"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/encryption"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/hashing"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/httpserver"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Run(config config.Config, logger logger.Logger) {
	postgresSQL := database.NewPostgreSQLDatabase(logger)

	encryptor := encryption.NewCryptoAES()
	hasher := hashing.New()

	casher := caching.NewRedis(caching.ConnectionOptions{
		Host:     config.Redis.Host,
		Password: config.Redis.Password,
		Database: config.Redis.Database,
	})

	serviceOptions := service.ServiceOptions{
		Logger:    logger,
		Config:    config,
		Cacher:    casher,
		Encryptor: encryptor,
		Hasher:    hasher,
		Database:  postgresSQL,
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
}
