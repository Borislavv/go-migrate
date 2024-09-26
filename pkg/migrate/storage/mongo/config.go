package mongo

import "github.com/kelseyhightower/envconfig"

type Configurator interface {
	GetMongoHost() string
	GetMongoPort() string
	GetMongoLogin() string
	GetMongoPassword() string
	GetMongoDatabase() string
	GetMongoMigrationsCollection() string
	GetMongoMigrationsDir() string
	IsMongoMigrationsEnabled() bool
}

type Config struct {
	MongoMigrationsEnabled    bool   `envconfig:"MONGO_MIGRATIONS_ENABLED" default:"false"`
	MongoHost                 string `envconfig:"MONGO_HOST"`
	MongoPort                 string `envconfig:"MONGO_PORT"`
	MongoLogin                string `envconfig:"MONGO_LOGIN"`
	MongoPassword             string `envconfig:"MONGO_PASSWORD"`
	MongoDatabase             string `envconfig:"MONGO_DATABASE"`
	MongoMigrationsCollection string `envconfig:"MONGO_MIGRATIONS_COLLECTION" default:"migrationVersions"`
	MongoMigrationsDir        string `envconfig:"MONGO_MIGRATIONS_DIR"`
}

func Load() (*Config, error) {
	cfg := new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) IsMongoMigrationsEnabled() bool {
	return c.MongoMigrationsEnabled
}

func (c *Config) GetMongoHost() string {
	return c.MongoHost
}

func (c *Config) GetMongoPort() string {
	return c.MongoPort
}

func (c *Config) GetMongoLogin() string {
	return c.MongoLogin
}

func (c *Config) GetMongoPassword() string {
	return c.MongoPassword
}

func (c *Config) GetMongoDatabase() string {
	return c.MongoDatabase
}

func (c *Config) GetMongoMigrationsCollection() string {
	return c.MongoMigrationsCollection
}

func (c *Config) GetMongoMigrationsDir() string {
	return c.MongoMigrationsDir
}
