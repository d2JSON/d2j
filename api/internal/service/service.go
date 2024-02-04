package service

import (
	"context"
	"errors"

	"github.com/VladPetriv/postgreSQL2JSON/config"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/caching"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/database"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/encryption"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/hashing"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/logger"
)

type Services struct {
	Database DatabaseService
}

type serviceContext struct {
	logger    logger.Logger
	config    config.Config
	cacher    caching.Cacher
	encryptor encryption.Encryptor
	hasher    hashing.Hasher
	database  database.Database
}

type ServiceOptions struct {
	Logger    logger.Logger
	Config    config.Config
	Cacher    caching.Cacher
	Encryptor encryption.Encryptor
	Hasher    hashing.Hasher
	Database  database.Database
}

type DatabaseService interface {
	TestDatabaseConnection(ctx context.Context, options DatabaseConnectionOptions) error
	ConnectToDatabase(ctx context.Context, options ConnectToDatabaseOptions) (string, error)
	ListDatabaseTables(ctx context.Context, options ListDatabaseTablesOptions) ([]string, error)
	ConvertDatabaseResultToJSON(ctx context.Context, options ConvertDatabaseResultToJSONOptions) (string, error)
}

type ConnectToDatabaseOptions struct {
	SecretKey                 string                    `json:"secretKey" binding:"required"`
	ConnectionSessionTime     string                    `json:"connectionSessionTime" binding:"required"`
	DatabaseConnectionOptions DatabaseConnectionOptions `json:"databaseConnectionOptions"`
}

type DatabaseConnectionOptions struct {
	Host           string `json:"host" binding:"required"`
	Port           int    `json:"port" binding:"required"`
	DatabaseName   string `json:"databaseName" binding:"required"`
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	SSLModeEnabled bool   `json:"sslModeEnabled"`
}

type ListDatabaseTablesOptions struct {
	SecretKey   string `json:"secretKey" binding:"required"`
	DatabaseKey string `json:"databaseKey" binding:"required"`
}

type ConvertDatabaseResultToJSONOptions struct {
	SecretKey   string `json:"secretKey" binding:"required"`
	DatabaseKey string `json:"databaseKey" binding:"required"`
	TableName   string `json:"tableName" binding:"required"`

	Fields []string `json:"fileds"`
	Limit  int      `json:"limit"`
	Where  string   `json:"where"`
}

var ErrConnectionSessionTimeExpired = errors.New("connection session time expired")
