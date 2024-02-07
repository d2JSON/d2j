package database

// Database ...
type Database interface {
	Connect(options ConnectionOptions) (DBClient, error)
	Close() error
}

// DBClient ...
type DBClient interface {
	ListTables() ([]Table, error)
	ExecuteQuery(query string) ([]string, error)
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

// Table represents a database table.
type Table struct {
	SchemaName string `db:"schemaname"`
	TableName  string `db:"tablename"`
}
