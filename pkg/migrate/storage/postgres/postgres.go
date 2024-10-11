package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/Borislavv/migrate/v4"
	"github.com/Borislavv/migrate/v4/database/postgres"
	"os"
)

const DriverName = "postgres"

type Postgres struct {
	ctx context.Context
	db  *sql.DB
	cfg Configurator
	fs  embed.FS
}

func New(ctx context.Context, cfg Configurator, fs embed.FS) (*Postgres, error) {
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

	return &Postgres{ctx: ctx, db: db, cfg: cfg, fs: fs}, nil
}

func (m *Postgres) Name() string {
	return DriverName
}

func (m *Postgres) Up() error {
	s, err := m.migrate()
	if err != nil {
		return err
	}

	if err = s.Up(); err != nil {
		return err
	}

	return nil
}

func (m *Postgres) Down() error {
	s, err := m.migrate()
	if err != nil {
		return err
	}

	if err = s.Down(); err != nil {
		return err
	}

	return nil
}

func (m *Postgres) migrate() (*migrate.Migrate, error) {
	if m.db == nil {
		return nil, errors.New("the underlying database pointer is not initialized, you need to call the 'New' method first")
	}

	d, err := postgres.WithInstance(m.ctx, m.db, &postgres.Config{
		DatabaseName:    m.cfg.GetPostgresDatabase(),
		MigrationsTable: m.cfg.GetPostgresMigrationsTable(),
	})
	if err != nil {
		return nil, err
	}

	if err = os.Mkdir(DriverName, 0777); err != nil {
		return nil, fmt.Errorf("could not create PostgreSQL migrations directory: %w", err)
	}

	if err = os.CopyFS(DriverName, m.fs); err != nil {
		return nil, fmt.Errorf("could not copy PostgreSQL migrations fs: %w", err)
	}

	s, err := migrate.NewWithDatabaseInstance("file://"+DriverName+"/migrations", DriverName, d)
	if err != nil {
		return nil, err
	}

	return s, nil
}
