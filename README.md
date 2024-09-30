## Migration library 

### Example: 
    func main() {
        if err := run(); err != nil {
            os.Exit(1)
        }
    }

    func run() error {
        ctx, cancel := context.WithCancel(context.Background())
        defer cancel()
    
        output, cancel, err := logger.NewOutput("")
        if err != nil {
            logger.JsonRawLog("logger: unable to initialize output", loggerenum.FatalLvl, err)
            return err
        }
        defer cancel()
    
        lgr, cancel, err := logger.NewLogrus(output)
        if err != nil {
            logger.JsonRawLog("logger: unable to initialize logrus", loggerenum.FatalLvl, err)
            return err
        }
        defer cancel()
    
        migrator, err := migrate.New(ctx, lgr, storage.NewFactory(lgr))
        if err != nil {
            return lgr.Fatal(ctx, errors.New("migrations: failed to init migrate"), logger.Fields{"err": err.Error()})
        }
    
        if err = migrator.Up(); err != nil {
            return lgr.Fatal(ctx, errors.New("migrations: up: completed with errors"), logger.Fields{"err": err.Error()})
        } else {
            lgr.InfoMsg(ctx, "migrations: up: completed", nil)
        }
    
        return nil
    }

### ENV:

#### MongoDB:
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
#### MySQL:
    type Config struct {
        MySQLMigrationsEnabled bool   `envconfig:"MYSQL_MIGRATIONS_ENABLED" default:"false"`
        MySQLHost              string `envconfig:"MYSQL_HOST"`
        MySQLPort              string `envconfig:"MYSQL_PORT"`
        MySQLUsername          string `envconfig:"MYSQL_LOGIN"`
        MySQLPassword          string `envconfig:"MYSQL_PASSWORD"`
        MySQLDatabase          string `envconfig:"MYSQL_DATABASE"`
        MySQLMigrationsTable   string `envconfig:"MYSQL_MIGRATIONS_TABLE" default:"migration_versions"`
        MySQLMigrationsDir     string `envconfig:"MYSQL_MIGRATIONS_DIR"`
    }
#### PostgreSQL:
    type Config struct {
        PostgresMigrationsEnabled bool   `envconfig:"POSTGRES_MIGRATIONS_ENABLED" default:"false"`
        PostgresHost              string `envconfig:"POSTGRES_HOST"`
        PostgresPort              string `envconfig:"POSTGRES_PORT"`
        PostgresUsername          string `envconfig:"POSTGRES_LOGIN"`
        PostgresPassword          string `envconfig:"POSTGRES_PASSWORD"`
        PostgresDatabase          string `envconfig:"POSTGRES_DATABASE"`
        PostgresMigrationsTable   string `envconfig:"POSTGRES_MIGRATIONS_TABLE" default:"migration_versions"`
        PostgresMigrationsDir     string `envconfig:"POSTGRES_MIGRATIONS_DIR"`
    }