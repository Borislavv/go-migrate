package mysql

import "github.com/kelseyhightower/envconfig"

type Configurator interface {
	IsMySQLMigrationsEnabled() bool
	GetMySQLUsername() string
	GetMySQLPassword() string
	GetMySQLDatabase() string
	GetMySQLHost() string
	GetMySQLPort() string
	GetMySQLMigrationsTable() string
}

type Config struct {
	MySQLMigrationsEnabled bool   `envconfig:"MYSQL_MIGRATIONS_ENABLED" default:"false"`
	MySQLHost              string `envconfig:"MYSQL_HOST"`
	MySQLPort              string `envconfig:"MYSQL_PORT"`
	MySQLUsername          string `envconfig:"MYSQL_LOGIN"`
	MySQLPassword          string `envconfig:"MYSQL_PASSWORD"`
	MySQLDatabase          string `envconfig:"MYSQL_DATABASE"`
	MySQLMigrationsTable   string `envconfig:"MYSQL_MIGRATIONS_TABLE" default:"migration_versions"`
}

func Load() (*Config, error) {
	cfg := new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) IsMySQLMigrationsEnabled() bool {
	return c.MySQLMigrationsEnabled
}

func (c *Config) GetMySQLHost() string {
	return c.MySQLHost
}

func (c *Config) GetMySQLPort() string {
	return c.MySQLPort
}

func (c *Config) GetMySQLUsername() string {
	return c.MySQLUsername
}

func (c *Config) GetMySQLPassword() string {
	return c.MySQLPassword
}

func (c *Config) GetMySQLDatabase() string {
	return c.MySQLDatabase
}

func (c *Config) GetMySQLMigrationsTable() string {
	return c.MySQLMigrationsTable
}
