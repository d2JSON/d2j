package database

import "github.com/VladPetriv/d2j/pkg/errs"

// Database ...
type Database interface {
	Connect(options ConnectionOptions) (DBClient, error)
	Close() error
}

// ConnectionOptions represents an options that used for connecting to the database.
type ConnectionOptions struct {
	Host           string
	Port           int
	Username       string
	Password       string
	DatabaseName   string
	SSLModeEnabled bool
}

// DBClient ...
type DBClient interface {
	BuildQuery(options BuildQueryOptions) string
	ListTables() ([]Table, error)
	ExecuteQuery(query string) ([]string, error)
}

// BuildQueryOptions represents an options that used for building query.
type BuildQueryOptions struct {
	TableName string
	Fields    []string
	Limit     int
	Where     string
}

// Table represents a database table.
type Table struct {
	SchemaName string `db:"schemaname"`
	TableName  string `db:"tablename"`
}

var (
	// ErrDatabaseDoesNotExists - database does not exists.
	ErrDatabaseDoesNotExists = errs.New("database does not exists")
	// ErrInvalidUsername - invalid username.
	ErrInvalidUsername = errs.New("invalid username")
	// ErrInvalidHost - invalid host.
	ErrInvalidHost = errs.New("invalid host")
	// ErrInvalidPort - invalid port.
	ErrInvalidPort = errs.New("invalid port")
	// ErrNoAccess - no access to database.
	ErrNoAccess = errs.New("no access")
)
