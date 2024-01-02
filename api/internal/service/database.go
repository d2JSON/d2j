package service

import (
	"context"
	"fmt"

	"github.com/VladPetriv/postgreSQL2JSON/pkg/database"
)

type databaseService struct {
	serviceContext
}

var _ DatabaseService = (*databaseService)(nil)

func NewDatabaseService(options *ServiceOptions) *databaseService {
	return &databaseService{
		serviceContext: serviceContext{
			logger:   options.Logger,
			config:   options.Config,
			database: options.Database,
		},
	}
}

func (d databaseService) TestDatabaseConnection(ctx context.Context, options TestDatabaseConnectionOptions) error {
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

	err = d.database.Close()
	if err != nil {
		logger.Error("close db connection", "err", err)
		return fmt.Errorf("close db connection: %w", err)
	}

	return nil
}
