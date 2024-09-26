package storage

import (
	"context"
	"errors"
	"github.com/Borislavv/go-logger/pkg/logger"
	"github.com/Borislavv/go-migrate/pkg/migrate/storage/mongo"
	"github.com/Borislavv/go-migrate/pkg/migrate/storage/mysql"
	"github.com/Borislavv/go-migrate/pkg/migrate/storage/postgres"
)

var (
	ErrFailedLoadMongoConfig    = errors.New("failed to load mongo config")
	ErrFailedLoadMySQLConfig    = errors.New("failed to load MySQL config")
	ErrFailedLoadPostgresConfig = errors.New("failed to load PostgreSQL config")
)

var (
	ErrFailedCreateInstanceMongo    = errors.New("failed to create mongo instance")
	ErrFailedCreateInstanceMySQL    = errors.New("failed to create MySQL instance")
	ErrFailedCreateInstancePostgres = errors.New("failed to create PostgreSQL instance")
)

var (
	ErrUnableToGetMongoMigrator    = errors.New("unable to get mongo migrator")
	ErrUnableToGetMySQLMigrator    = errors.New("unable to get MySQL migrator")
	ErrUnableToGetPostgresMigrator = errors.New("unable to get PostgreSQL migrator")
)

var (
	ErrMongoMigrationsIsNotEnabled    = errors.New("mongo migrations is not enabled")
	ErrMySQLMigrationsIsNotEnabled    = errors.New("MySQL migrations is not enabled")
	ErrPostgresMigrationsIsNotEnabled = errors.New("PostgreSQL migrations is not enabled")
)

type Factorier interface {
	Make(ctx context.Context) ([]Storager, error)
}

type Factory struct {
	logger logger.Logger
}

func NewFactory(logger logger.Logger) *Factory {
	return &Factory{logger: logger}
}

func (f *Factory) Make(ctx context.Context) ([]Storager, error) {
	storages := make([]Storager, 0, 3)

	if s, err := f.getMongo(ctx); err != nil {
		if !errors.Is(err, ErrMongoMigrationsIsNotEnabled) {
			return nil, f.logger.Error(ctx, ErrUnableToGetMongoMigrator, logger.Fields{
				"err": err.Error(),
			})
		}
	} else {
		storages = append(storages, s)
	}

	if s, err := f.getMySQL(ctx); err != nil {
		if !errors.Is(err, ErrMySQLMigrationsIsNotEnabled) {
			return nil, f.logger.Error(ctx, ErrUnableToGetMySQLMigrator, logger.Fields{
				"err": err.Error(),
			})
		}
	} else {
		storages = append(storages, s)
	}

	if s, err := f.getPostgres(ctx); err != nil {
		if !errors.Is(err, ErrPostgresMigrationsIsNotEnabled) {
			return nil, f.logger.Error(ctx, ErrUnableToGetPostgresMigrator, logger.Fields{
				"err": err.Error(),
			})
		}
	} else {
		storages = append(storages, s)
	}

	return storages, nil
}

func (f *Factory) getMongo(ctx context.Context) (Storager, error) {
	c, err := mongo.Load()
	if err != nil {
		return nil, f.logger.Error(ctx, ErrFailedLoadMongoConfig, logger.Fields{
			"err": err.Error(),
		})
	}
	if c.MongoMigrationsEnabled {
		m, err := mongo.New(ctx, c)
		if err != nil {
			return nil, f.logger.Error(ctx, ErrFailedCreateInstanceMongo, logger.Fields{
				"err": err.Error(),
			})
		}
		return m, nil
	}

	return nil, ErrMongoMigrationsIsNotEnabled
}

func (f *Factory) getMySQL(ctx context.Context) (Storager, error) {
	c, err := mysql.Load()
	if err != nil {
		return nil, f.logger.Error(ctx, ErrFailedLoadMySQLConfig, logger.Fields{
			"err": err.Error(),
		})
	}
	if c.MySQLMigrationsEnabled {
		m, err := mysql.New(ctx, c)
		if err != nil {
			return nil, f.logger.Error(ctx, ErrFailedCreateInstanceMySQL, logger.Fields{
				"err": err.Error(),
			})
		}
		return m, nil
	}

	return nil, ErrMySQLMigrationsIsNotEnabled
}

func (f *Factory) getPostgres(ctx context.Context) (Storager, error) {
	c, err := postgres.Load()
	if err != nil {
		return nil, f.logger.Error(ctx, ErrFailedLoadPostgresConfig, logger.Fields{
			"err": err.Error(),
		})
	}
	if c.PostgresMigrationsEnabled {
		m, err := postgres.New(ctx, c)
		if err != nil {
			return nil, f.logger.Error(ctx, ErrFailedCreateInstancePostgres, logger.Fields{
				"err": err.Error(),
			})
		}
		return m, nil
	}

	return nil, ErrPostgresMigrationsIsNotEnabled
}
