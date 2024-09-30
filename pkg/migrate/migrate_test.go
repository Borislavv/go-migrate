package migrate

import (
	"context"
	"github.com/Borislavv/go-logger/pkg/logger"
	loggerenum "github.com/Borislavv/go-logger/pkg/logger/enum"
	"github.com/Borislavv/go-migrate/pkg/migrate/storage"
	"testing"
)

type TestFactory struct {
}

func (f *TestFactory) Make(_ context.Context) ([]storage.Storager, error) {
	return []storage.Storager{&TestStorage{}}, nil
}

type TestStorage struct {
}

func (s *TestStorage) Name() string {
	return "test"
}
func (s *TestStorage) Up() error {
	return nil
}
func (s *TestStorage) Down() error {
	return nil
}

func TestMigrate_Up(t *testing.T) {
	out, cancel, err := logger.NewOutput(loggerenum.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	lgr, cancel, err := logger.NewLogrus(out)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	m, err := New(context.Background(), lgr, &TestFactory{})
	if err != nil {
		t.Fatal(err)
	}

	if err = m.Up(); err != nil {
		t.Fatal(err)
	}
}
