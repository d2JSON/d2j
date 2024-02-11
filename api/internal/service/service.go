package service

import (
	"context"
	"errors"

	"github.com/VladPetriv/d2j/config"
	"github.com/VladPetriv/d2j/pkg/caching"
	"github.com/VladPetriv/d2j/pkg/database"
	"github.com/VladPetriv/d2j/pkg/encryption"
	"github.com/VladPetriv/d2j/pkg/hashing"
	"github.com/VladPetriv/d2j/pkg/logger"
)

// Services represents a structure that contains all application services.
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

// Options represents a structure that contains all packages that needed for services.
type Options struct {
	Logger    logger.Logger
	Config    config.Config
	Cacher    caching.Cacher
	Encryptor encryption.Encryptor
	Hasher    hashing.Hasher
	Database  database.Database
}

// DatabaseService ...
type DatabaseService interface {
	TestDatabaseConnection(ctx context.Context, options DatabaseConnectionOptions) error
	ConnectToDatabase(ctx context.Context, options ConnectToDatabaseOptions) (string, error)
	ListDatabaseTables(ctx context.Context, options ListDatabaseTablesOptions) ([]string, error)
	ConvertDatabaseResultToJSON(ctx context.Context, options ConvertDatabaseResultToJSONOptions) (string, error)
}

// ConnectToDatabaseOptions represents options for ConnectToDatabase method.
type ConnectToDatabaseOptions struct {
	SecretKey                 string                    `json:"secretKey" binding:"required"`
	ConnectionSessionTime     string                    `json:"connectionSessionTime" binding:"required"`
	DatabaseConnectionOptions DatabaseConnectionOptions `json:"databaseConnectionOptions"`
}

// DatabaseConnectionOptions represents options that required for creating connection with database.
type DatabaseConnectionOptions struct {
	Host           string `json:"host" binding:"required"`
	Port           int    `json:"port" binding:"required"`
	DatabaseName   string `json:"databaseName" binding:"required"`
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	SSLModeEnabled bool   `json:"sslModeEnabled"`
}

// ListDatabaseTablesOptions represents options for ListDatabaseTables method.
type ListDatabaseTablesOptions struct {
	SecretKey   string `json:"secretKey" binding:"required"`
	DatabaseKey string `json:"databaseKey" binding:"required"`
}

// ConvertDatabaseResultToJSONOptions represents options for ConvertDatabaseResultToJS method.
type ConvertDatabaseResultToJSONOptions struct {
	SecretKey   string `json:"secretKey" binding:"required"`
	DatabaseKey string `json:"databaseKey" binding:"required"`
	TableName   string `json:"tableName" binding:"required"`

	Fields []string `json:"fileds"`
	Limit  int      `json:"limit"`
	Where  string   `json:"where"`
}

// ErrConnectionSessionTimeExpired occurs when entered by user time for connection session is expired.
var ErrConnectionSessionTimeExpired = errors.New("connection session time expired")
