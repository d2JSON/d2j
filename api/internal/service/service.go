package service

import (
	"context"

	"github.com/VladPetriv/postgreSQL2JSON/config"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/database"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/logger"
)

type Services struct {
	Database DatabaseService
}

type serviceContext struct {
	logger   logger.Logger
	config   config.Config
	database database.Database
}

type ServiceOptions struct {
	Logger   logger.Logger
	Config   config.Config
	Database database.Database
}

type DatabaseService interface {
	TestDatabaseConnection(ctx context.Context, options TestDatabaseConnectionOptions) error
}

type TestDatabaseConnectionOptions struct {
	Host           string `json:"host" binding:"required"`
	Port           int    `json:"port" binding:"required"`
	DatabaseName   string `json:"databaseName" binding:"required"`
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	SSLModeEnabled bool   `json:"sslModeEnabled"`
}
