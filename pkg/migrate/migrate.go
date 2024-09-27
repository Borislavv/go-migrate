package migrate

import (
	"context"
	"errors"
	"github.com/Borislavv/go-logger/pkg/logger"
	"github.com/Borislavv/go-migrate/pkg/migrate/storage"
	"github.com/golang-migrate/migrate/v4"
	"golang.org/x/sync/errgroup"
)

var (
	ErrNoChanges       = migrate.ErrNoChange
	ErrMigratorFactory = errors.New("failed to make migrators")
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

	return &Migrate{
		ctx:      ctx,
		logger:   logger,
		storages: storages,
	}, nil
}

// Up executes each migrator in parallel wrapping them in errgroup without context.
//   - withCtx determines whether the errgroup will be spawned with the context (if so, the context will be terminated
//     when the first error occurs, and other migrators that did not have an error will be closed).
func (m *Migrate) Up(withCtx bool) error {
	eg := &errgroup.Group{}
	ctx := context.Background()
	if withCtx {
		eg, ctx = errgroup.WithContext(m.ctx)
	}

	for _, migrator := range m.storages {
		eg.Go(func() error {
			prefix := "migrations: " + migrator.Name() + ": up: "

			if err := migrator.Up(); err != nil {
				if errors.Is(err, ErrNoChanges) {
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
//   - withCtx determines whether the errgroup will be spawned with the context (if so, the context will be terminated
//     when the first error occurs, and other migrators that did not have an error will be closed).
func (m *Migrate) Down(withCtx bool) error {
	eg := &errgroup.Group{}
	ctx := context.Background()
	if withCtx {
		eg, ctx = errgroup.WithContext(m.ctx)
	}

	for _, migrator := range m.storages {
		eg.Go(func() error {
			prefix := "migrations: " + migrator.Name() + ": down: "

			if err := migrator.Up(); err != nil {
				if errors.Is(err, ErrNoChanges) {
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
