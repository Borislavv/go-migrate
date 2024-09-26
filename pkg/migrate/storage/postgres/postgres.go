package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Borislavv/migrate/v4"
	"github.com/Borislavv/migrate/v4/database/postgres"
	"os"
	"path/filepath"
)

const DriverName = "postgres"

type Postgres struct {
	ctx context.Context
	db  *sql.DB
	cfg Configurator
}

func New(ctx context.Context, cfg Configurator) (*Postgres, error) {
	dsn := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		DriverName,
		cfg.GetPostgresUsername(),
		cfg.GetPostgresPassword(),
		cfg.GetPostgresHost(),
		cfg.GetPostgresPort(),
		cfg.GetPostgresDatabase(),
	)

	db, err := sql.Open(DriverName, dsn)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		_ = db.Close()
	}()

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Postgres{ctx: ctx, db: db, cfg: cfg}, nil
}

func (m *Postgres) Name() string {
	return DriverName
}

func (m *Postgres) Up() error {
	if m.db == nil {
		return errors.New("the underlying database pointer is not initialized, you need to call the 'New' method first")
	}

	d, err := postgres.WithInstance(m.ctx, m.db, &postgres.Config{
		DatabaseName:    m.cfg.GetPostgresDatabase(),
		MigrationsTable: m.cfg.GetPostgresMigrationsTable(),
	})
	if err != nil {
		return err
	}

	r, err := os.Getwd()
	if err != nil {
		return err
	}

	s, err := migrate.NewWithDatabaseInstance("file://"+filepath.Join(r, m.cfg.GetPostgresMigrationsDir()), DriverName, d)
	if err != nil {
		return err
	}

	if err = s.Up(); err != nil {
		return err
	}

	return nil
}

func (m *Postgres) Down() error {
	if m.db == nil {
		return errors.New("the underlying database pointer is not initialized, you need to call the 'New' method first")
	}

	d, err := postgres.WithInstance(m.ctx, m.db, &postgres.Config{
		DatabaseName:    m.cfg.GetPostgresDatabase(),
		MigrationsTable: m.cfg.GetPostgresMigrationsTable(),
	})
	if err != nil {
		return err
	}

	r, err := os.Getwd()
	if err != nil {
		return err
	}

	s, err := migrate.NewWithDatabaseInstance("file://"+filepath.Join(r, m.cfg.GetPostgresMigrationsDir()), DriverName, d)
	if err != nil {
		return err
	}

	if err = s.Down(); err != nil {
		return err
	}

	return nil
}
