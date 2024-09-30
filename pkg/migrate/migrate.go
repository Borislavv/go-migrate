package migrate

import (
	"context"
	"errors"
	"github.com/Borislavv/go-logger/pkg/logger"
	"github.com/Borislavv/go-migrate/pkg/migrate/storage"
	"github.com/Borislavv/migrate/v4"
	"golang.org/x/sync/errgroup"
)

var (
	ErrMigratorFactory         = errors.New("failed to make migrators")
	ErrNoOneMigratorWasDefined = errors.New("no migrators were defined")
)

type Migrate struct {
	ctx      context.Context
	logger   logger.Logger
	storages []storage.Storager
}

func New(ctx context.Context, logger logger.Logger, factory storage.Factorier) (*Migrate, error) {
	storages, err := factory.Make(ctx)
	if err != nil {
		return nil, logger.Fatal(ctx, ErrMigratorFactory, nil)
	}

	if len(storages) == 0 {
		return nil, logger.Fatal(ctx, ErrNoOneMigratorWasDefined, nil)
	}

	return &Migrate{
		ctx:      ctx,
		logger:   logger,
		storages: storages,
	}, nil
}

// Up executes each migrator in parallel wrapping them in errgroup without context.
func (m *Migrate) Up() error {
	eg := &errgroup.Group{}
	ctx := context.Background()

	for _, migrator := range m.storages {
		eg.Go(func() error {
			prefix := "migrations: " + migrator.Name() + ": up: "

			if err := migrator.Up(); err != nil {
				if errors.Is(err, migrate.ErrNoChange) {
					m.logger.InfoMsg(ctx, prefix+"no changes detected", logger.Fields{
						"storage": migrator.Name(),
					})
					return nil
				}

				return m.logger.Fatal(ctx, errors.New(prefix+"error occurred while applying migrations"), logger.Fields{
					"err":     err.Error(),
					"storage": migrator.Name(),
				})
			}

			m.logger.InfoMsg(ctx, prefix+"schema successfully upped", logger.Fields{
				"storage": migrator.Name(),
			})
			return nil
		})
	}

	return eg.Wait()
}

// Down executes each migrator in parallel wrapping them in errgroup without context.
func (m *Migrate) Down() error {
	eg := &errgroup.Group{}
	ctx := context.Background()

	for _, migrator := range m.storages {
		eg.Go(func() error {
			prefix := "migrations: " + migrator.Name() + ": down: "

			if err := migrator.Down(); err != nil {
				if errors.Is(err, migrate.ErrNoChange) {
					m.logger.InfoMsg(ctx, prefix+"no changes detected", logger.Fields{
						"storage": migrator.Name(),
					})
					return nil
				}

				return m.logger.Fatal(ctx, errors.New(prefix+"error occurred while applying migrations"), logger.Fields{
					"err":     err.Error(),
					"storage": migrator.Name(),
				})
			}

			m.logger.InfoMsg(ctx, prefix+"schema successfully downgraded", logger.Fields{
				"storage": migrator.Name(),
			})
			return nil
		})
	}

	return eg.Wait()
}

// Migrators returns all for self management.
func (m *Migrate) Migrators() []storage.Storager {
	return m.storages
}
