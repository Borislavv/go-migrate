package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/Borislavv/migrate/v4"
	"github.com/Borislavv/migrate/v4/database/mongodb"
	_ "github.com/Borislavv/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"path/filepath"
)

const DriverName = "mongodb"

type Mongo struct {
	db  *mongo.Database
	cfg Configurator
}

func New(ctx context.Context, cfg Configurator) (*Mongo, error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf(
		"%s://%s:%s@%s:%s/?authSource=%s",
		DriverName,
		cfg.GetMongoLogin(),
		cfg.GetMongoPassword(),
		cfg.GetMongoHost(),
		cfg.GetMongoPort(),
		cfg.GetMongoDatabase(),
	))

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		_ = mongoClient.Disconnect(ctx)
	}()

	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &Mongo{db: mongoClient.Database(cfg.GetMongoDatabase()), cfg: cfg}, nil
}

func (m *Mongo) Name() string {
	return DriverName
}

func (m *Mongo) Up() error {
	if m.db == nil {
		return errors.New("the underlying database pointer is not initialized, you need to call the 'New' method first")
	}

	d, err := mongodb.WithInstance(m.db.Client(), &mongodb.Config{
		DatabaseName:         m.db.Name(),
		MigrationsCollection: m.cfg.GetMongoMigrationsCollection(),
	})
	if err != nil {
		return err
	}

	r, err := os.Getwd()
	if err != nil {
		return err
	}

	s, err := migrate.NewWithDatabaseInstance("file://"+filepath.Join(r, m.cfg.GetMongoMigrationsDir()), DriverName, d)
	if err != nil {
		return err
	}

	if err = s.Up(); err != nil {
		return err
	}

	return nil
}

func (m *Mongo) Down() error {
	if m.db == nil {
		return errors.New("the underlying database pointer is not initialized, you need to call the 'New' method first")
	}

	d, err := mongodb.WithInstance(m.db.Client(), &mongodb.Config{
		DatabaseName:         m.db.Name(),
		MigrationsCollection: m.cfg.GetMongoMigrationsDir(),
	})
	if err != nil {
		return err
	}

	r, err := os.Getwd()
	if err != nil {
		return err
	}

	s, err := migrate.NewWithDatabaseInstance("file://"+filepath.Join(r, m.cfg.GetMongoMigrationsDir()), DriverName, d)
	if err != nil {
		return err
	}

	if err = s.Down(); err != nil {
		return err
	}

	return nil
}
