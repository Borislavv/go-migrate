package postgres

import "github.com/kelseyhightower/envconfig"

type Configurator interface {
	IsPostgresMigrationsEnabled() bool
	GetPostgresUsername() string
	GetPostgresPassword() string
	GetPostgresDatabase() string
	GetPostgresHost() string
	GetPostgresPort() string
	GetPostgresMigrationsTable() string
}

type Config struct {
	PostgresMigrationsEnabled bool   `envconfig:"POSTGRES_MIGRATIONS_ENABLED" default:"false"`
	PostgresHost              string `envconfig:"POSTGRES_HOST"`
	PostgresPort              string `envconfig:"POSTGRES_PORT"`
	PostgresUsername          string `envconfig:"POSTGRES_LOGIN"`
	PostgresPassword          string `envconfig:"POSTGRES_PASSWORD"`
	PostgresDatabase          string `envconfig:"POSTGRES_DATABASE"`
	PostgresMigrationsTable   string `envconfig:"POSTGRES_MIGRATIONS_TABLE" default:"migration_versions"`
}

func Load() (*Config, error) {
	cfg := new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) IsPostgresMigrationsEnabled() bool {
	return c.PostgresMigrationsEnabled
}

func (c *Config) GetPostgresHost() string {
	return c.PostgresHost
}

func (c *Config) GetPostgresPort() string {
	return c.PostgresPort
}

func (c *Config) GetPostgresUsername() string {
	return c.PostgresUsername
}

func (c *Config) GetPostgresPassword() string {
	return c.PostgresPassword
}

func (c *Config) GetPostgresDatabase() string {
	return c.PostgresDatabase
}

func (c *Config) GetPostgresMigrationsTable() string {
	return c.PostgresMigrationsTable
}
