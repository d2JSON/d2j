package controller

import (
	"fmt"
	"net/http"

	"github.com/VladPetriv/d2j/config"
	"github.com/VladPetriv/d2j/internal/service"
	"github.com/VladPetriv/d2j/pkg/logger"
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

// New creates and register all routes.
func New(options Options) {
	options.Handler.Use(corsMiddleware)

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
				if options.Config.HTTP.SendDetailsOnInternalError {
					c.AbortWithStatusJSON(http.StatusInternalServerError, err)
				} else {
					//  Do not send error details to client
					err := c.AbortWithError(http.StatusInternalServerError, err)
					if err != nil {
						logger.Error("failed to abort with error", "err", err)
					}
				}
			} else {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
			}
			return
		}

		c.JSON(http.StatusOK, body)
	}
}

func corsMiddleware(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")

	c.Header("Access-Control-Allow-Origin", origin)
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
