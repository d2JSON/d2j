package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VladPetriv/postgreSQL2JSON/config"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/httpserver"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Run(config config.Config, logger logger.Logger) {
	httpHandler := gin.New()

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
		log.Println("app - Run - signal: " + s.String())

	case err := <-httpServer.Notify():
		log.Println("app - Run - httpServer.Notify", "err", err)
	}
}
