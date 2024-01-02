package database

type Database interface {
	Connect(options ConnectionOptions) (DBClient, error)
	Close() error
}

type DBClient interface{}

type ConnectionOptions struct {
	Host           string
	Port           int
	Username       string
	Password       string
	DatabaseName   string
	SSLModeEnabled bool
}

type Table struct {
	Name string `json:"name"`
}
