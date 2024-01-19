package controller

import (
	"github.com/VladPetriv/postgreSQL2JSON/internal/service"
	"github.com/gin-gonic/gin"
)

type databaseRouter struct {
	RouterContext
}

func setupDatabaseRoutes(options RouterOptions) {
	r := databaseRouter{
		RouterContext: RouterContext{
			logger:   options.Logger,
			config:   options.Config,
			services: options.Services,
		},
	}

	dbGroup := options.Handler.Group("/db")
	{
		dbGroup.POST("/test-connection", wrapHandler(options, r.testDBConnection))
		dbGroup.POST("/connect", wrapHandler(options, r.connectToDatabase))
	}
}

type testDBConnectionRequestBody struct {
	*service.DatabaseConnectionOptions
}

func (r databaseRouter) testDBConnection(c *gin.Context) (interface{}, *httpResponseError) {
	logger := r.logger.Named("databaseRouter.TestDBConnection")

	var requestBody testDBConnectionRequestBody
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		logger.Error("bind request body to json", "err", err)
		return nil, &httpResponseError{Message: "invalid request body", Type: ErrorTypeClient}
	}
	logger.Debug("parsed request body", "requestBody", requestBody)

	err = r.services.Database.TestDatabaseConnection(c, *requestBody.DatabaseConnectionOptions)
	if err != nil {
		logger.Error("test database connection", "err", err)
		return nil, &httpResponseError{Message: "test database connection failed", Type: ErrorTypeServer}
	}

	logger.Info("connection tested")
	return nil, nil
}

type connectToDatabaseRequestBody struct {
	*service.ConnectToDatabaseOptions
}

type connectToDatabaseResponse struct {
	SecretKey string `json:"secretKey"`
}

func (r databaseRouter) connectToDatabase(c *gin.Context) (interface{}, *httpResponseError) {
	logger := r.logger.Named("databaseRouter.connectToDatabase")

	var requestBody connectToDatabaseRequestBody
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		logger.Error("bind request body to json", "err", err)
		return nil, &httpResponseError{Message: "invalid request body", Type: ErrorTypeClient}
	}
	logger.Debug("parsed request body", "requestBody", requestBody)

	key, err := r.services.Database.ConnectToDatabase(c, *requestBody.ConnectToDatabaseOptions)
	if err != nil {
		logger.Error("connect to database", "err", err)
		return nil, &httpResponseError{Message: "connect to database", Type: ErrorTypeServer}
	}

	logger.Info("connected to database", "key", key)
	return connectToDatabaseResponse{SecretKey: key}, nil
}
