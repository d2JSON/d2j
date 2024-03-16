package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/VladPetriv/d2j/pkg/caching"
	"github.com/VladPetriv/d2j/pkg/database"
	"github.com/VladPetriv/d2j/pkg/encryption"
	"github.com/VladPetriv/d2j/pkg/errs"
	"github.com/google/uuid"
)

type databaseService struct {
	serviceContext
}

var _ DatabaseService = (*databaseService)(nil)

// NewDatabaseService  is used to create an instance of database service.
func NewDatabaseService(options *Options) *databaseService {
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
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return d.handleConnectionErrors(err)
		}

		logger.Error("connect to database", "err", err)
		return fmt.Errorf("connect to database: %w", err)
	}
	logger.Debug("connected to database")

	err = d.database.Close()
	if err != nil {
		logger.Error("close database connection", "err", err)
		return fmt.Errorf("close database connection: %w", err)
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
		logger.Error("connect to database", "err", err)
		return "", fmt.Errorf("connect to database: %w", err)
	}
	logger.Debug("connected to database")

	defer func() {
		err = d.database.Close()
		if err != nil {
			logger.Error("close database connection", "err", err)
		} else {
			logger.Debug("closed database connection")
		}
	}()

	databaseKey, err := d.hasher.GenerateHashFromString(uuid.NewString())
	if err != nil {
		logger.Debug("generate hash from string", "err", err)
		return "", fmt.Errorf("generate hash from string: %w", err)
	}
	logger.Debug("generated hash")

	marshalledConnectionOptions, err := json.Marshal(options.DatabaseConnectionOptions)
	if err != nil {
		logger.Error("marshal connection options to JSON", "err", err)
		return "", fmt.Errorf("marshal connection options to JSON: %w", err)
	}
	logger.Debug("marshalled connection options to JSON")

	encryptedConnectionData, err := d.encryptor.Encrypt(encryption.EncryptOptions{
		Data:   marshalledConnectionOptions,
		Secret: d.config.App.EncryptionSecretKey,
	})
	if err != nil {
		logger.Error("encrypt connection options", "err", err)
		return "", fmt.Errorf("encrypt connection options: %w", err)
	}
	logger.Debug("encrypted connection data")

	connectionSessionTime, err := time.ParseDuration(options.ConnectionSessionTime)
	if err != nil {
		logger.Error("parse connection session time", "err", err)
		return "", fmt.Errorf("parse connection session time: %w", err)
	}
	logger.Debug("parsed connection session time")

	err = d.cacher.Write(ctx, caching.WriteOptions{
		Key:   databaseKey,
		Value: encryptedConnectionData,
		TTL:   connectionSessionTime,
	})
	if err != nil {
		logger.Error("write connection data to cache", "err", err)
		return "", fmt.Errorf("write connection data to cache: %w", err)
	}
	logger.Debug("wrote connection data to cache")

	return databaseKey, nil
}

func (d databaseService) ListDatabaseTables(ctx context.Context, options ListDatabaseTablesOptions) ([]string, error) {
	logger := d.logger.Named("databaseService.ListDatabaseTables")

	databaseConnectionOptions, err := d.getDatabaseCredentials(ctx, options.DatabaseKey)
	if err != nil {
		if errors.Is(err, ErrConnectionSessionTimeExpired) {
			logger.Info("connection session time expired")
			return nil, err
		}
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
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, d.handleConnectionErrors(err)
		}

		logger.Error("connect to database", "err", err)
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	logger.Debug("connected to database")

	databaseTables, err := databaseClient.ListTables()
	if err != nil {
		logger.Error("list database tables", "err", err)
		return nil, fmt.Errorf("list database tables: %w", err)
	}
	logger.Debug("got database tables", "databaseTables", databaseTables)

	tableNames := make([]string, len(databaseTables))

	for i := range tableNames {
		tableNames[i] = databaseTables[i].TableName
	}
	logger.Debug("converted database tables to slice of strings", "tableNames", tableNames)

	return tableNames, nil
}

func (d databaseService) ConvertDatabaseResultToJSON(ctx context.Context, options ConvertDatabaseResultToJSONOptions) (string, error) {
	logger := d.logger.Named("databaseService.ConvertDatabaseResultToJSON")

	databaseConnectionOptions, err := d.getDatabaseCredentials(ctx, options.DatabaseKey)
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
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return "", d.handleConnectionErrors(err)
		}
		logger.Error("connect to database", "err", err)
		return "", fmt.Errorf("connect to database: %w", err)
	}
	logger.Debug("connected to database")

	query := databaseClient.BuildQuery(database.BuildQueryOptions{
		TableName: options.TableName,
		Limit:     options.Limit,
		Fields:    options.Fields,
		Where:     options.Where,
	})
	logger.Debug("built query", "query", query)

	databaseResult, err := databaseClient.ExecuteQuery(query)
	if err != nil {
		logger.Error("execute query", "err", err)
		return "", fmt.Errorf("execute query: %w", err)
	}
	logger.Debug("got database result", "databaseResult", databaseResult)

	JSONResult := "[ "

	for index, row := range databaseResult {
		// Do not add comma for the last slice element to get valid JSON format.
		if index == len(databaseResult)-1 {
			JSONResult += fmt.Sprintf("%s\n", row)

			continue
		}

		JSONResult += fmt.Sprintf("%s,\n", row)
	}

	JSONResult += " ]"
	logger.Debug("converted database result to JSON", "JSONResult", JSONResult)

	return JSONResult, nil
}

func (d databaseService) handleConnectionErrors(err error) error {
	if errs.HasAnyGivenMessage(err, database.ErrDatabaseDoesNotExists.Error()) {
		return ErrDatabaseDoesNotExists
	}
	if errs.HasAnyGivenMessage(err, database.ErrInvalidUsername.Error()) {
		return ErrInvalidUsername
	}
	if errs.HasAnyGivenMessage(err, database.ErrInvalidHost.Error()) {
		return ErrInvalidHost
	}
	if errs.HasAnyGivenMessage(err, database.ErrInvalidPort.Error()) {
		return ErrInvalidPort
	}
	if errs.HasAnyGivenMessage(err, database.ErrNoAccess.Error()) {
		return ErrNoAccessToDatabase
	}

	return err
}

func (d databaseService) getDatabaseCredentials(ctx context.Context, databaseKey string) (*DatabaseConnectionOptions, error) {
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
		Secret:        d.config.App.EncryptionSecretKey,
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
