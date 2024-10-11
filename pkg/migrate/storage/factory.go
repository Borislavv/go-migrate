package storage

import (
	"context"
	"embed"
	"errors"
	"github.com/Borislavv/go-logger/pkg/logger"
	"github.com/Borislavv/go-migrate/pkg/migrate/storage/mongo"
	"github.com/Borislavv/go-migrate/pkg/migrate/storage/mysql"
	"github.com/Borislavv/go-migrate/pkg/migrate/storage/postgres"
)

var (
	ErrFailedLoadMongoDBConfig  = errors.New("failed to load MongoDB config")
	ErrFailedLoadMySQLConfig    = errors.New("failed to load MySQL config")
	ErrFailedLoadPostgresConfig = errors.New("failed to load PostgreSQL config")
)

var (
	ErrFailedCreateInstanceMongoDB  = errors.New("failed to create MongoDB instance")
	ErrFailedCreateInstanceMySQL    = errors.New("failed to create MySQL instance")
	ErrFailedCreateInstancePostgres = errors.New("failed to create PostgreSQL instance")
)

var (
	ErrUnableToGetMongoDBMigrator  = errors.New("unable to get MongoDB migrator")
	ErrUnableToGetMySQLMigrator    = errors.New("unable to get MySQL migrator")
	ErrUnableToGetPostgresMigrator = errors.New("unable to get PostgreSQL migrator")
)

var (
	ErrMongoDBMigrationsIsNotEnabled  = errors.New("MongoDB migrations is not enabled")
	ErrMySQLMigrationsIsNotEnabled    = errors.New("MySQL migrations is not enabled")
	ErrPostgresMigrationsIsNotEnabled = errors.New("PostgreSQL migrations is not enabled")
)

var (
	ErrMongoDBFSWasOmitted    = errors.New("unable to migrate MongoDB, target filesystem was omitted")
	ErrMySQLFSWasOmitted      = errors.New("unable to migrate MySQL, target filesystem was omitted")
	ErrPostgreSQLFSWasOmitted = errors.New("unable to migrate PostgreSQL, target filesystem was omitted")
)

type Storage int

const (
	MongoDB Storage = iota
	MySQL
	PostgreSQL
)

type Filesystems map[Storage]embed.FS

type Factory struct {
	logger      logger.Logger
	filesystems Filesystems
}

func NewFactory(logger logger.Logger, filesystems Filesystems) *Factory {
	return &Factory{logger: logger, filesystems: filesystems}
}

func (f *Factory) Make(ctx context.Context) ([]Storager, error) {
	storages := make([]Storager, 0, 3)

	if s, err := f.getMongo(ctx); err != nil {
		if !errors.Is(err, ErrMongoDBMigrationsIsNotEnabled) {
			return nil, f.logger.Error(ctx, ErrUnableToGetMongoDBMigrator, logger.Fields{
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
	cfg, err := mongo.Load()
	if err != nil {
		return nil, f.logger.Error(ctx, ErrFailedLoadMongoDBConfig, logger.Fields{
			"err": err.Error(),
		})
	}
	if cfg.MongoMigrationsEnabled {
		fs, ok := f.filesystems[MongoDB]
		if !ok {
			return nil, f.logger.Error(ctx, ErrMongoDBFSWasOmitted, nil)
		}

		m, err := mongo.New(ctx, cfg, fs)
		if err != nil {
			return nil, f.logger.Error(ctx, ErrFailedCreateInstanceMongoDB, logger.Fields{
				"err": err.Error(),
			})
		}
		return m, nil
	}

	return nil, ErrMongoDBMigrationsIsNotEnabled
}

func (f *Factory) getMySQL(ctx context.Context) (Storager, error) {
	cfg, err := mysql.Load()
	if err != nil {
		return nil, f.logger.Error(ctx, ErrFailedLoadMySQLConfig, logger.Fields{
			"err": err.Error(),
		})
	}
	if cfg.MySQLMigrationsEnabled {
		fs, ok := f.filesystems[MySQL]
		if !ok {
			return nil, f.logger.Error(ctx, ErrMySQLFSWasOmitted, nil)
		}

		m, err := mysql.New(ctx, cfg, fs)
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
	cfg, err := postgres.Load()
	if err != nil {
		return nil, f.logger.Error(ctx, ErrFailedLoadPostgresConfig, logger.Fields{
			"err": err.Error(),
		})
	}
	if cfg.PostgresMigrationsEnabled {
		fs, ok := f.filesystems[PostgreSQL]
		if !ok {
			return nil, f.logger.Error(ctx, ErrPostgreSQLFSWasOmitted, nil)
		}

		m, err := postgres.New(ctx, cfg, fs)
		if err != nil {
			return nil, f.logger.Error(ctx, ErrFailedCreateInstancePostgres, logger.Fields{
				"err": err.Error(),
			})
		}
		return m, nil
	}

	return nil, ErrPostgresMigrationsIsNotEnabled
}
