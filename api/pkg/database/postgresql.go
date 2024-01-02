package database

import (
	"fmt"

	"github.com/VladPetriv/postgreSQL2JSON/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type postgreSQL struct {
	logger logger.Logger

	db *sqlx.DB
}

var _ Database = (*postgreSQL)(nil)

func NewPostgreSQLDatabase(logger logger.Logger) *postgreSQL {
	return &postgreSQL{
		logger: logger,
	}
}

type postgreSQLClient struct {
	logger logger.Logger

	db *sqlx.DB
}

var _ DBClient = (*postgreSQLClient)(nil)

func newPostgreSQLClient(logger logger.Logger, db *sqlx.DB) *postgreSQLClient {
	return &postgreSQLClient{
		logger: logger,
		db:     db,
	}
}

func (p *postgreSQL) Connect(options ConnectionOptions) (DBClient, error) {
	logger := p.logger.Named("postgreSQL.Connect")

	connectionString := buildConnectionString(options)
	logger.Debug("built connection string", "connectionString", connectionString)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		logger.Error("connect to postgresql", "err", err)
		return nil, fmt.Errorf("connect to postgresql: %w", err)
	}

	p.db = db

	logger.Info("connected to postgresql")
	return newPostgreSQLClient(logger, db), nil
}

func buildConnectionString(options ConnectionOptions) string {
	connectionString := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d",
		options.Username, options.Password, options.DatabaseName, options.Host, options.Port,
	)

	if !options.SSLModeEnabled {
		connectionString += " sslmode=disable"
	}
	if options.SSLModeEnabled {
		connectionString += " sslmode=enable"
	}

	return connectionString
}

func (p *postgreSQL) Close() error {
	logger := p.logger.Named("postgreSQL.Close")

	err := p.db.DB.Close()
	if err != nil {
		logger.Error("close postgresql connection", err)
		return fmt.Errorf("close postgresql connection: %w", err)
	}

	logger.Info("closed postgresql connection")
	return nil
}
