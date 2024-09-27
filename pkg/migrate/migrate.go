package migrate

import (
	"context"
	"errors"
	"github.com/Borislavv/go-logger/pkg/logger"
	loggerenum "github.com/Borislavv/go-logger/pkg/logger/enum"
	"github.com/Borislavv/go-migrate/pkg/migrate/storage"
	"github.com/golang-migrate/migrate/v4"
	"golang.org/x/sync/errgroup"
)

var (
	ErrNoChanges       = migrate.ErrNoChange
	ErrMigratorFactory = errors.New("failed to make migrators")
)

type log struct {
	msg    string
	level  string
	fields logger.Fields
}

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
func (m *Migrate) Up() error {
	eg := &errgroup.Group{}
	ctx := context.Background()

	logsCh := make(chan *log, len(m.storages))
	defer close(logsCh)

	go func() {
		for l := range logsCh {
			m.logger.LogMsg(ctx, l.msg, l.level, l.fields)
		}
	}()

	for _, migrator := range m.storages {
		eg.Go(func() error {
			prefix := "migrations: " + migrator.Name() + ": up: "

			if err := migrator.Up(); err != nil {
				if errors.Is(err, ErrNoChanges) {
					logsCh <- &log{
						msg:   prefix + "no changes detected",
						level: loggerenum.InfoLvl,
						fields: logger.Fields{
							"storage": migrator.Name(),
						},
					}
					return nil
				}

				err = errors.New(prefix + "error occurred while applying migrations")
				logsCh <- &log{
					msg:   err.Error(),
					level: loggerenum.FatalLvl,
					fields: logger.Fields{
						"err":     err.Error(),
						"storage": migrator.Name(),
					},
				}
				return err
			}

			logsCh <- &log{
				msg:   prefix + "schema successfully upped",
				level: loggerenum.InfoLvl,
				fields: logger.Fields{
					"storage": migrator.Name(),
				},
			}
			return nil
		})
	}

	return eg.Wait()
}

// Down executes each migrator in parallel wrapping them in errgroup without context.
func (m *Migrate) Down() error {
	eg := &errgroup.Group{}
	ctx := context.Background()
	logsCh := make(chan *log)

	for _, migrator := range m.storages {
		eg.Go(func() error {
			prefix := "migrations: " + migrator.Name() + ": down: "

			if err := migrator.Up(); err != nil {
				if errors.Is(err, ErrNoChanges) {
					logsCh <- &log{
						msg:   prefix + "no changes detected",
						level: loggerenum.InfoLvl,
						fields: logger.Fields{
							"storage": migrator.Name(),
						},
					}
					return nil
				}

				err = errors.New(prefix + "error occurred while applying migrations")
				logsCh <- &log{
					msg:   err.Error(),
					level: loggerenum.FatalLvl,
					fields: logger.Fields{
						"err":     err.Error(),
						"storage": migrator.Name(),
					},
				}
				return err
			}

			logsCh <- &log{
				msg:   prefix + "schema successfully downgraded",
				level: loggerenum.InfoLvl,
				fields: logger.Fields{
					"storage": migrator.Name(),
				},
			}
			return nil
		})
	}

	resultCh := make(chan error, 1)
	defer close(resultCh)

	go func() {
		resultCh <- eg.Wait()
		close(logsCh)
	}()

	for l := range logsCh {
		m.logger.LogMsg(ctx, l.msg, l.level, l.fields)
	}

	return <-resultCh
}

// Migrators returns all for self management.
func (m *Migrate) Migrators() []storage.Storager {
	return m.storages
}
