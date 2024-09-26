package migrate

import (
	"context"
	"errors"
	"github.com/Borislavv/go-logger/pkg/logger"
	"github.com/Borislavv/go-migrate/pkg/migrate/storage"
	"sync"
)

var (
	ErrMigratorFactory = errors.New("failed to make migrators")
	ErrMigrationFailed = errors.New("migration failed")
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

func (m *Migrate) Up(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for _, migrator := range m.storages {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := migrator.Up(); err != nil {
				m.logger.ErrorMsg(ctx, migrator.Name()+": up: "+ErrMigrationFailed.Error(), logger.Fields{
					"err": err.Error(),
				})
			}
		}()
	}

	return nil
}

func (m *Migrate) Down(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for _, migrator := range m.storages {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := migrator.Down(); err != nil {
				m.logger.ErrorMsg(ctx, migrator.Name()+": down: "+ErrMigrationFailed.Error(), logger.Fields{
					"err": err.Error(),
				})
			}
		}()
	}

	return nil
}
