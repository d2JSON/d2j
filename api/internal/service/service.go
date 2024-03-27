package service

import (
	"context"
	"errors"

	"github.com/VladPetriv/d2j/config"
	"github.com/VladPetriv/d2j/pkg/caching"
	"github.com/VladPetriv/d2j/pkg/database"
	"github.com/VladPetriv/d2j/pkg/encryption"
	"github.com/VladPetriv/d2j/pkg/errs"
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
	DatabaseKey string `json:"databaseKey" binding:"required"`
}

// ConvertDatabaseResultToJSONOptions represents options for ConvertDatabaseResultToJS method.
type ConvertDatabaseResultToJSONOptions struct {
	DatabaseKey string `json:"databaseKey" binding:"required"`
	TableName   string `json:"tableName" binding:"required"`

	Fields []string `json:"fields"`
	Limit  int      `json:"limit"`
	Where  string   `json:"where"`
}

var (
	// ErrConnectionSessionTimeExpired occurs when entered by user time for connection session is expired.
	ErrConnectionSessionTimeExpired = errors.New("connection session time expired")
	// ErrDatabaseDoesNotExists occurs when entered database name does not exists.
	ErrDatabaseDoesNotExists = errs.New("The entered database does not exist. Please verify the database name and try again")
	// ErrInvalidUsername occurs when user enters incorrect username for database.
	ErrInvalidUsername = errs.New("Invalid username. Please check and try again.")
	// ErrInvalidHost occurs when user enters incorrect database host.
	ErrInvalidHost = errs.New("Invalid database host. Please check and try again.")
	// ErrInvalidPort occurs when user enters incorrect database port.
	ErrInvalidPort = errs.New("Invalid port number. Please check and try again.")
	// ErrNoAccessToDatabase occurs when user database config does not allow connection from different IPs.
	ErrNoAccessToDatabase = errs.New("Access denied. Please verify your credentials and connection settings")
)
