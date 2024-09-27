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

type log struct {
	msg    string
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
//   - withCtx determines whether the errgroup will be spawned with the context (if so, the context will be terminated
//     when the first error occurs, and other migrators that did not have an error will be closed).
func (m *Migrate) Up(withCtx bool) error {
	eg := &errgroup.Group{}
	ctx := context.Background()
	if withCtx {
		eg, ctx = errgroup.WithContext(m.ctx)
	}

	infoLogsCh := make(chan *log)
	defer close(infoLogsCh)
	go func() {
		for l := range infoLogsCh {
			m.logger.InfoMsg(ctx, l.msg, l.fields)
		}
	}()

	fatalLogsCh := make(chan *log)
	defer close(fatalLogsCh)
	go func() {
		for l := range fatalLogsCh {
			m.logger.FatalMsg(ctx, l.msg, l.fields)
		}
	}()

	for _, migrator := range m.storages {
		eg.Go(func() error {
			prefix := "migrations: " + migrator.Name() + ": up: "

			if err := migrator.Up(); err != nil {
				if errors.Is(err, ErrNoChanges) {
					infoLogsCh <- &log{
						msg: prefix + "no changes detected",
						fields: logger.Fields{
							"storage": migrator.Name(),
						},
					}
					return nil
				}

				err = errors.New(prefix + "error occurred while applying migrations")
				fatalLogsCh <- &log{
					msg: err.Error(),
					fields: logger.Fields{
						"err":     err.Error(),
						"storage": migrator.Name(),
					},
				}
				return err
			}

			infoLogsCh <- &log{
				msg: prefix + "schema successfully upped",
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
//   - withCtx determines whether the errgroup will be spawned with the context (if so, the context will be terminated
//     when the first error occurs, and other migrators that did not have an error will be closed).
func (m *Migrate) Down(withCtx bool) error {
	eg := &errgroup.Group{}
	ctx := context.Background()
	if withCtx {
		eg, ctx = errgroup.WithContext(m.ctx)
	}

	infoLogsCh := make(chan *log)
	defer close(infoLogsCh)
	go func() {
		for l := range infoLogsCh {
			m.logger.InfoMsg(ctx, l.msg, l.fields)
		}
	}()

	fatalLogsCh := make(chan *log)
	defer close(fatalLogsCh)
	go func() {
		for l := range fatalLogsCh {
			m.logger.FatalMsg(ctx, l.msg, l.fields)
		}
	}()

	for _, migrator := range m.storages {
		eg.Go(func() error {
			prefix := "migrations: " + migrator.Name() + ": down: "

			if err := migrator.Up(); err != nil {
				if errors.Is(err, ErrNoChanges) {
					infoLogsCh <- &log{
						msg: prefix + "no changes detected",
						fields: logger.Fields{
							"storage": migrator.Name(),
						},
					}
					return nil
				}

				err = errors.New(prefix + "error occurred while applying migrations")
				fatalLogsCh <- &log{
					msg: err.Error(),
					fields: logger.Fields{
						"err":     err.Error(),
						"storage": migrator.Name(),
					},
				}
				return err
			}

			infoLogsCh <- &log{
				msg: prefix + "schema successfully downgraded",
				fields: logger.Fields{
					"storage": migrator.Name(),
				},
			}
			return nil
		})
	}

	return eg.Wait()
}

// Migrators returns all for self management.
func (m *Migrate) Migrators() []storage.Storager {
	return m.storages
}
