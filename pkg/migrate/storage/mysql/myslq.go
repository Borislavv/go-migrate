package mysql

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/Borislavv/migrate/v4"
	"github.com/Borislavv/migrate/v4/database/mysql"
	"os"
	"path/filepath"
)

const DriverName = "mysql"

type MySQL struct {
	ctx context.Context
	db  *sql.DB
	cfg Configurator
	fs  embed.FS
}

func New(ctx context.Context, cfg Configurator, fs embed.FS) (*MySQL, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		cfg.GetMySQLUsername(),
		cfg.GetMySQLPassword(),
		cfg.GetMySQLHost(),
		cfg.GetMySQLPort(),
		cfg.GetMySQLDatabase(),
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

	return &MySQL{ctx: ctx, db: db, cfg: cfg, fs: fs}, nil
}

func (m *MySQL) Name() string {
	return DriverName
}

func (m *MySQL) Up() error {
	s, err := m.migrate()
	if err != nil {
		return err
	}

	if err = s.Up(); err != nil {
		return err
	}

	return nil
}

func (m *MySQL) Down() error {
	s, err := m.migrate()
	if err != nil {
		return err
	}

	if err = s.Down(); err != nil {
		return err
	}

	return nil
}

func (m *MySQL) migrate() (*migrate.Migrate, error) {
	if m.db == nil {
		return nil, errors.New("the underlying database pointer is not initialized, you need to call the 'New' method first")
	}

	d, err := mysql.WithInstance(m.ctx, m.db, &mysql.Config{
		DatabaseName:    m.cfg.GetMySQLDatabase(),
		MigrationsTable: m.cfg.GetMySQLMigrationsTable(),
	})
	if err != nil {
		return nil, err
	}

	rootDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	destDir := filepath.Join(rootDir, "tmp", DriverName)
	if err = os.MkdirAll(destDir, 0777); err != nil {
		return nil, fmt.Errorf("could not create MySQL migrations directory: %w", err)
	}

	if err = os.CopyFS(destDir, m.fs); err != nil {
		return nil, fmt.Errorf("could not copy MySQL migrations fs: %w", err)
	}

	migrationsDir := filepath.Join(destDir, "migrations")
	s, err := migrate.NewWithDatabaseInstance("file://"+migrationsDir, DriverName, d)
	if err != nil {
		return nil, err
	}

	return s, nil
}
