package database

type Database interface {
	Connect(options ConnectionOptions) (DBClient, error)
	Close() error
}

type DBClient interface {
	ListTables() ([]Table, error)
	ExecuteQuery(query string) ([]string, error)
}

type ConnectionOptions struct {
	Host           string
	Port           int
	Username       string
	Password       string
	DatabaseName   string
	SSLModeEnabled bool
}

type Table struct {
	SchemaName string `db:"schemaname"`
	TableName  string `db:"tablename"`
}
