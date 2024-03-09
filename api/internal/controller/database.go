package controller

import (
	"errors"

	"github.com/VladPetriv/d2j/internal/service"
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

	database := options.Handler.Group("/database")
	{
		database.POST("/test-connection", wrapHandler(options, r.testDBConnection))
		database.POST("/connect", wrapHandler(options, r.connectToDatabase))
		database.POST("/list-tables", wrapHandler(options, r.listDatabaseTables))
		database.POST("/get-json", wrapHandler(options, r.convertDatabaseResultToJSON))
	}
}

type testDBConnectionRequestBody struct {
	*service.DatabaseConnectionOptions
}

type testDBConnectionResponseBody struct {
	Message string `json:"message"`
}

func (r databaseRouter) testDBConnection(c *gin.Context) (interface{}, *httpResponseError) {
	logger := r.logger.Named("databaseRouter.testDBConnection")

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
		return nil, &httpResponseError{Message: "Failed to test connection with your database!\nPlease try again!", Type: ErrorTypeServer}
	}

	logger.Info("connection tested")
	return testDBConnectionResponseBody{
		Message: "We have successfully established a connection with your database.",
	}, nil
}

type connectToDatabaseRequestBody struct {
	*service.ConnectToDatabaseOptions
}

type connectToDatabaseResponse struct {
	DatabaseKey string `json:"databaseKey"`
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

	databaseKey, err := r.services.Database.ConnectToDatabase(c, *requestBody.ConnectToDatabaseOptions)
	if err != nil {
		logger.Error("connect to database", "err", err)
		return nil, &httpResponseError{Message: "connect to database", Type: ErrorTypeServer}
	}

	logger.Info("connected to database", "databaseKey", databaseKey)
	return connectToDatabaseResponse{databaseKey}, nil
}

type listDatabaseTablesRequestBody struct {
	*service.ListDatabaseTablesOptions
}

type listDatabaseTablesResponse struct {
	Tables []string `json:"tables"`
}

func (r databaseRouter) listDatabaseTables(c *gin.Context) (interface{}, *httpResponseError) {
	logger := r.logger.Named("databaseRouter.listDatabaseTables")

	var requestBody listDatabaseTablesRequestBody
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		logger.Error("bind request body to json", "err", err)
		return nil, &httpResponseError{Message: "invalid request body", Type: ErrorTypeClient}
	}
	logger.Debug("parsed request body", "requestBody", requestBody)

	tables, err := r.services.Database.ListDatabaseTables(c, *requestBody.ListDatabaseTablesOptions)
	if err != nil {
		if errors.Is(err, service.ErrConnectionSessionTimeExpired) {
			logger.Info(err.Error())
			return nil, &httpResponseError{Message: "Connection session time expired", Type: ErrorTypeClient}
		}

		logger.Error("list database tables", "err", err)
		return nil, &httpResponseError{Message: "list database tables", Type: ErrorTypeServer}
	}

	logger.Info("got database tables", "tables", tables)
	return listDatabaseTablesResponse{tables}, nil
}

type convertDatabaseResultToJSONRequestBody struct {
	*service.ConvertDatabaseResultToJSONOptions
}

type convertDatabaseResultToJSONResponseBody struct {
	Result string `json:"result"`
}

func (r databaseRouter) convertDatabaseResultToJSON(c *gin.Context) (interface{}, *httpResponseError) {
	logger := r.logger.Named("databaseRouter.convertDatabaseResultToJSON")

	var reqBody convertDatabaseResultToJSONRequestBody
	err := c.ShouldBindJSON(&reqBody)
	if err != nil {
		logger.Error("bind request body to json", "err", err)
		return nil, &httpResponseError{Message: "invalid request body", Type: ErrorTypeClient}
	}
	logger.Debug("parsed request body", "reqBody", reqBody)

	convertedResult, err := r.services.Database.ConvertDatabaseResultToJSON(c, *reqBody.ConvertDatabaseResultToJSONOptions)
	if err != nil {
		if errors.Is(err, service.ErrConnectionSessionTimeExpired) {
			logger.Info(err.Error())
			return nil, &httpResponseError{Message: "Connection session time expired", Type: ErrorTypeClient}
		}

		logger.Error("convert database result to JSON", "err", err)
		return nil, &httpResponseError{Message: "convert database result to JSON", Type: ErrorTypeServer}
	}

	logger.Info("converted database result to JSON")
	return convertDatabaseResultToJSONResponseBody{convertedResult}, nil
}
