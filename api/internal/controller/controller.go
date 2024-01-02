package controller

import (
	"fmt"
	"net/http"

	"github.com/VladPetriv/postgreSQL2JSON/config"
	"github.com/VladPetriv/postgreSQL2JSON/internal/service"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Options is used to parameterize http controller via NewHTTPController.
type Options struct {
	Handler  *gin.Engine
	Logger   logger.Logger
	Config   config.Config
	Services service.Services
}

// RouterContext provides a shared context for all routers.
type RouterContext struct {
	logger   logger.Logger
	config   config.Config
	services service.Services
}

// RouterOptions provides shared options for all routers.
type RouterOptions struct {
	Handler  *gin.RouterGroup
	Logger   logger.Logger
	Config   config.Config
	Services service.Services
}

func New(options Options) {
	routerOptions := RouterOptions{
		Handler:  options.Handler.Group("/api"),
		Logger:   options.Logger,
		Config:   options.Config,
		Services: options.Services,
	}

	// Routers
	{
		setupDatabaseRoutes(routerOptions)
	}
}

// httpResponseError provides a base error type for all errors.
type httpResponseError struct {
	Type    httpErrType `json:"-"`
	Message string      `json:"message"`
}

// httpErrType is used to define error type.
type httpErrType string

const (
	// ErrorTypeServer is an "unexpected" internal server error.
	ErrorTypeServer httpErrType = "server"
	// ErrorTypeClient is an "expected" business error.
	ErrorTypeClient httpErrType = "client"
)

// Error is used to convert an error to a string.
func (e httpResponseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// wrapHandler provides unified error handling for all handlers.
func wrapHandler(options RouterOptions, handler func(c *gin.Context) (interface{}, *httpResponseError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := options.Logger.Named("controller.wrapHandler")

		// execute handler
		body, err := handler(c)

		// check if middleware
		if body == nil && err == nil {
			return
		}

		// check error
		if err != nil {
			if err.Type == ErrorTypeServer {
				err := c.AbortWithError(http.StatusInternalServerError, err)
				if err != nil {
					logger.Error("failed to abort with error", "err", err)
				}
			} else {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
			}
			return
		}

		c.JSON(http.StatusOK, body)
	}
}
