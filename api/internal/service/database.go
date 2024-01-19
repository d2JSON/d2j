package service

import (
	"context"
	"encoding/json"
	"fmt"

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

	err = d.cacher.Write(ctx, secretKey, encryptedConnectionOPtions)
	if err != nil {
		logger.Error("write connection options to cache", "err", err)
		return "", fmt.Errorf("write connection options to cache: %w", err)
	}
	logger.Debug("wrote connection options to cache")

	return secretKey, nil
}
