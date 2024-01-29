package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/VladPetriv/postgreSQL2JSON/pkg/caching"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/database"
	"github.com/VladPetriv/postgreSQL2JSON/pkg/encryption"
	"github.com/google/uuid"
)

type databaseService struct {
	serviceContext
}

var _ DatabaseService = (*databaseService)(nil)

func NewDatabaseService(options *ServiceOptions) *databaseService {
	return &databaseService{
		serviceContext: serviceContext{
			logger:    options.Logger,
			config:    options.Config,
			database:  options.Database,
			cacher:    options.Cacher,
			hasher:    options.Hasher,
			encryptor: options.Encryptor,
		},
	}
}

func (d databaseService) TestDatabaseConnection(ctx context.Context, options DatabaseConnectionOptions) error {
	logger := d.logger.Named("databaseService.TestDatabaseConnection")

	_, err := d.database.Connect(database.ConnectionOptions{
		Host:           options.Host,
		Port:           options.Port,
		Username:       options.Username,
		Password:       options.Password,
		DatabaseName:   options.DatabaseName,
		SSLModeEnabled: options.SSLModeEnabled,
	})
	if err != nil {
		logger.Error("connect to db", "err", err)
		return fmt.Errorf("connect to db: %w", err)
	}
	logger.Debug("connected to database")

	err = d.database.Close()
	if err != nil {
		logger.Error("close db connection", "err", err)
		return fmt.Errorf("close db connection: %w", err)
	}
	logger.Debug("closed database connection")

	return nil
}

func (d databaseService) ConnectToDatabase(ctx context.Context, options ConnectToDatabaseOptions) (string, error) {
	logger := d.logger.Named("databaseService.ConnectToDatabase")

	_, err := d.database.Connect(database.ConnectionOptions{
		Host:           options.DatabaseConnectionOptions.Host,
		Port:           options.DatabaseConnectionOptions.Port,
		Username:       options.DatabaseConnectionOptions.Username,
		Password:       options.DatabaseConnectionOptions.Password,
		DatabaseName:   options.DatabaseConnectionOptions.DatabaseName,
		SSLModeEnabled: options.DatabaseConnectionOptions.SSLModeEnabled,
	})
	if err != nil {
		logger.Error("connect to db", "err", err)
		return "", fmt.Errorf("connect to db: %w", err)
	}
	logger.Debug("connected to database")

	defer func() {
		err = d.database.Close()
		if err != nil {
			logger.Error("close db connection", "err", err)
		}
		logger.Debug("closed database connection")
	}()

	secretKey, err := d.hasher.GenerateHashFromString(uuid.NewString())
	if err != nil {
		logger.Debug("generate hash from string", "err", err)
		return "", fmt.Errorf("generate hash from string: %w", err)
	}
	logger.Debug("hash generated")

	marshalledConnectionOptions, err := json.Marshal(options.DatabaseConnectionOptions)
	if err != nil {
		logger.Error("marshal connection options to JSON", "err", err)
		return "", fmt.Errorf("marshal connection options to JSON: %w", err)
	}
	logger.Debug("marshalled connection options to JSON")

	encryptedConnectionOPtions, err := d.encryptor.Encrypt(encryption.EncryptOptions{
		Data:   marshalledConnectionOptions,
		Secret: options.SecretKey,
	})
	if err != nil {
		logger.Error("encrypt connection options", "err", err)
		return "", fmt.Errorf("encrypt connection options: %w", err)
	}

	ttl, err := time.ParseDuration(options.ConnectionSessionTime)
	if err != nil {
		logger.Error("parse connection session time", "err", err)
		return "", fmt.Errorf("parsr connection session time: %w", err)
	}

	err = d.cacher.Write(ctx, caching.WriteOptions{
		Key:   secretKey,
		Value: encryptedConnectionOPtions,
		TTL:   ttl,
	})
	if err != nil {
		logger.Error("write connection options to cache", "err", err)
		return "", fmt.Errorf("write connection options to cache: %w", err)
	}
	logger.Debug("wrote connection options to cache")

	return secretKey, nil
}

func (d databaseService) ListDatabaseTables(ctx context.Context, options ListDatabaseTablesOptions) ([]string, error) {
	logger := d.logger.Named("databaseService.ListDatabaseTables")

	databaseConnectionOptions, err := d.getDatabaseCredentials(ctx, options.SecretKey, options.DatabaseKey)
	if err != nil {
		logger.Error("get database credentials", "err", err)
		return nil, fmt.Errorf("get database credentials: %w", err)
	}
	logger.Debug("got database credentials")

	databaseClient, err := d.database.Connect(database.ConnectionOptions{
		Host:           databaseConnectionOptions.Host,
		Port:           databaseConnectionOptions.Port,
		Username:       databaseConnectionOptions.Username,
		Password:       databaseConnectionOptions.Password,
		DatabaseName:   databaseConnectionOptions.DatabaseName,
		SSLModeEnabled: databaseConnectionOptions.SSLModeEnabled,
	})
	if err != nil {
		logger.Error("connect to database", "err", err)
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	logger.Debug("connected to database")

	databaseTables, err := databaseClient.ListTables()
	if err != nil {
		logger.Error("list database tables", "err", err)
		return nil, fmt.Errorf("list database tables: %w", err)
	}

	tableNames := make([]string, len(databaseTables))

	for i := range tableNames {
		tableNames[i] = databaseTables[i].TableName
	}
	logger.Debug("converted database tables to slice of strings", "tableNames", tableNames)

	return tableNames, nil
}

func (d databaseService) ConvertDatabaseResultToJSON(ctx context.Context, options ConvertDatabaseResultToJSONOptions) (string, error) {
	logger := d.logger.Named("databaseService.ConvertDatabaseResultToJSON")

	databaseConnectionOptions, err := d.getDatabaseCredentials(ctx, options.SecretKey, options.DatabaseKey)
	if err != nil {
		logger.Error("get database credentials", "err", err)
		return "", fmt.Errorf("get database credentials: %w", err)
	}
	logger.Debug("got database credentials")

	databaseClient, err := d.database.Connect(database.ConnectionOptions{
		Host:           databaseConnectionOptions.Host,
		Port:           databaseConnectionOptions.Port,
		Username:       databaseConnectionOptions.Username,
		Password:       databaseConnectionOptions.Password,
		DatabaseName:   databaseConnectionOptions.DatabaseName,
		SSLModeEnabled: databaseConnectionOptions.SSLModeEnabled,
	})
	if err != nil {
		logger.Error("connect to database", "err", err)
		return "", fmt.Errorf("connect to database: %w", err)
	}
	logger.Debug("connected to database")

	query := d.buildQuery(buildQueryOptions{
		TableName: options.TableName,
		Limit:     options.Limit,
		Fields:    options.Fields,
		Where:     options.Where,
	})
	logger.Debug("built query", "query", query)

	queryResult, err := databaseClient.ExecuteQuery(query)
	if err != nil {
		logger.Error("execute query", "err", err)
		return "", fmt.Errorf("execute query: %w", err)
	}
	logger.Debug("got result after query", "queryResult", queryResult)

	JSONResult := "[ "

	for index, row := range queryResult {
		// Do not add comma for the last slice element to get valid JSON format.
		if index == len(queryResult)-1 {
			JSONResult += fmt.Sprintf("%s\n", row)

			continue
		}

		JSONResult += fmt.Sprintf("%s,\n", row)
	}

	JSONResult += " ]"
	logger.Debug("converted database query result to JSON", "JSONResult", JSONResult)

	formattedJSON, err := json.MarshalIndent(JSONResult, "", "\t")
	if err != nil {
		logger.Error("marshal intent JSON", "err", err)
		return "", fmt.Errorf("marshal intent JSON: %w", err)
	}
	logger.Debug("formatted JSON", "formattedJSON", string(formattedJSON))

	return string(formattedJSON), nil
}

type buildQueryOptions struct {
	TableName string
	Fields    []string
	Limit     int
	Where     string
}

func (d databaseService) buildQuery(options buildQueryOptions) string {
	logger := d.logger.Named("databaseService.buildQuery")
	var query string

	if len(options.Fields) == 0 {
		query = fmt.Sprintf("SELECT to_jsonb(%s) FROM %s", options.TableName, options.TableName)
		logger.Debug("built select and from statement", "query", query)
	}

	if len(options.Fields) != 0 {
		query += "SELECT jsonb_agg(jsonb_build_object("

		for i, f := range options.Fields {
			if i == len(options.Fields)-1 {
				query += fmt.Sprintf("'%s', %s))", f, f)

				continue
			}

			query += fmt.Sprintf("'%s', %s, ", f, f)
		}

		query += fmt.Sprintf("FROM %s", options.TableName)
		logger.Debug("built selece with specific fields and from statement", "query", query)
	}

	if len(options.Where) != 0 {
		query += fmt.Sprintf(" WHERE %s", options.Where)
		logger.Debug("added where condition", "query", query)
	}
	if options.Limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", options.Limit)
		logger.Debug("added limit", "query", query)
	}
	logger.Debug("built query", "query", query)

	return query
}

func (d databaseService) getDatabaseCredentials(ctx context.Context, secretKey, databaseKey string) (*DatabaseConnectionOptions, error) {
	logger := d.logger.Named("databaseService.getDatabaseCredentials")

	encryptedDatabaseCredentials, err := d.cacher.Read(ctx, databaseKey)
	if err != nil {
		if errors.Is(err, caching.ErrResultIsNil) {
			logger.Info("connection session time expired")
			return nil, ErrConnectionSessionTimeExpired
		}

		logger.Error("read encrypted database credentials from cache", "err", err)
		return nil, fmt.Errorf("read encrypted database credentials from cache: %w", err)
	}
	logger.Debug("read encrypted database credentials from cache")

	decryptedDatabaseCredentials, err := d.encryptor.Decrypt(encryption.DecryptOptions{
		EncryptedData: encryptedDatabaseCredentials,
		Secret:        secretKey,
	})
	if err != nil {
		logger.Error("decrypt database credentials", "err", err)
		return nil, fmt.Errorf("decrypt database credentials: %w", err)
	}
	logger.Debug("decrypted database credentials")

	var databaseConnectionOptions DatabaseConnectionOptions
	err = json.Unmarshal(decryptedDatabaseCredentials, &databaseConnectionOptions)
	if err != nil {
		logger.Error("unmarshal database credentials", "err", err)
		return nil, fmt.Errorf("unmarshal database credentials: %w", err)

	}
	logger.Debug("unmarshalled database credentials")

	return &databaseConnectionOptions, nil
}
