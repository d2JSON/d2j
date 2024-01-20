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

	encryptedDatabaseCredentials, err := d.cacher.Read(ctx, options.DatabaseKey)
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
		Secret:        options.SecretKey,
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
