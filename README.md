# go-pgsql
[![PkgGoDev](https://godoc.org/github.com/i4erkasov/go-pgsql?status.svg)](https://pkg.go.dev/github.com/i4erkasov/go-pgsql?tab=doc)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/i4erkasov/go-pgsql)
[![Go Report Card](https://goreportcard.com/badge/github.com/i4erkasov/go-pgsql)](https://goreportcard.com/report/github.com/i4erkasov/go-pgsql)

`go-pgsql` is a comprehensive Go package that facilitates managing PostgreSQL connection pools, database migrations, and transaction management. It leverages the robust `pgx` library for pool management and `golang-migrate` for database migrations, providing a full suite of tools for working with PostgreSQL in Go applications.

## Features

- *[Connection Pool Management](#connection-pool-management)*: Simplifies the creation and management of connection pools using `pgxpool`.
- *[Transaction Management](#transaction-management)*: Offers a convenient transaction manager that supports nested transactions and automatic error handling.
- *[Database Migrations](#migrate)*: Allows for smooth database schema migrations using the `golang-migrate` package.

## External Packages

This package uses the following external libraries:

- Connection Pooling: [github.com/jackc/pgx/v4/pgxpool](https://github.com/jackc/pgx)
- Database Migrations: [github.com/golang-migrate/migrate/v4](https://github.com/golang-migrate/migrate)
- Configuration: [github.com/spf13/viper](https://github.com/spf13/viper)

## Installation

```bash
go get -u github.com/i4erkasov/go-pgsql
```

## Connection Pool Management

#### Using Viper for Configuration

```yaml
pgsql:
  pgpool:
    default:
      nodes: [
        "postgres://user:password@127.0.0.1:5432/db?sslmode=disable&client_encoding=UTF8"
        "postgres://user:password@127.0.1.2:5432/db?sslmode=disable&client_encoding=UTF8"
        "postgres://user:password@127.1.1.3:5432/db?sslmode=disable&client_encoding=UTF8"
      ]
      max_conns: 5 # Default 4
      min_conns: 1 # Default 0
      max_conn_lifetime: "90m" # Default: 1 hour
      max_conn_idle_time: "31m" # Default: 30 minutes
      health_check_period: "1m" # Default: 1 minute
      lazy_conn: false # Optional
      prefer_simple_protocol: true # Optional
```

```go
import (
    "github.com/spf13/viper"
    "github.com/i4erkasov/go-pgsql/pgxpool"
)

func main() {
    viper.SetConfigName("config")
    viper.AddConfigPath(".")
    viper.AutomaticEnv()

    err := viper.ReadInConfig()
    if err != nil {
        // Handle configuration read error
    }

    registry, err := pgxpool.NewWithViper(viper.GetViper())
    if err != nil {
        // Handle registry creation error
    }

    // Use the registry
}
```

#### Using Direct Configurations

```go
import "github.com/i4erkasov/go-pgsql/pgxpool"

func main() {
    registry, err := pgxpool.NewWithConfigOptions(
        WithConfig(pgxpool.DEFAULT, Config{/* ... configuration params ... */}),
        WithConfig("pool1", Config{/* ... configuration params ... */}),
        WithConfig("pool2", Config{/* ... configuration params ... */}),
        // Other configuration...
    )

    if err != nil {
        // Handle registry creation error
    }
	
    //Use the registry
}
```
#### or

```go
import "github.com/i4erkasov/go-pgsql/pgxpool"

func main() {
    config := pgxpool.Config{
    // Set your configuration here
    }

    registry, err := pgxpool.NewWithConfigs(map[string]pgxpool.Config{
        "default": config,
    })
    if err != nil {
        // Handle registry creation error
    }

    // Use the registry
}
```

Configuration Parameters
- `nodes`: (Required): Array of database connection strings.
- `max_conns`: Maximum number of connections in the pool (default: 4).
- `min_conns`: Minimum number of connections in the pool (default: 0).
- `max_conn_lifetime`: Maximum lifetime of a connection (default: 1 hour).
- `max_conn_idle_time`: Maximum idle time of a connection (default: 30 minutes).
- `health_check_period`: Frequency of health checks for idle connections (default: 1 minute).
- `lazy_conn`: Whether to establish a new connection lazily (default: false).
- `prefer_simple_protoco`l: Whether to use simple protocol for new connections (default: false).


```go
func main() {
    // ...
	
    // Use the registry
	
    pool, err := registry.Pools() // is DEFAULT pool getter
    // or
    pool1, err := registry.GetPoolName("pool1") // is pool getter by name.
    
    master := pool.Maser() // returns master connections pool
    slave := pool.Slave() // returns slave connections pool
}
```

## Transaction Management

The package includes a sophisticated transaction manager that allows for simple and complex transactional operations, including support for nested transactions.

### Using

#### Starting Transactions

To begin a transaction, simply call the `Begin` method on an instance of `TxManager`:

```go
import (
    "github.com/i4erkasov/go-pgsql/tx"
)

func main() {
    txManager, err := tx.NewTxManager(registry)
    if err != nil {
    // Handle error
    }

    // Use the txManager
}
```

#### Beginning Transaction 

```go
txCtx, err := txManager.Begin(context.Background())
if err != nil {
    // Handle error
}
```

#### Committing Transactions
Once your operations are complete, `Commit` the transaction:

```go
err = txManager.Commit(txCtx)
if err != nil {
    // Handle error
}
```

#### Rolling Back Transactions
If an error occurs, or you need to abort the transaction, call `Rollback`:

```go
err = txManager.Rollback(txCtx)
if err != nil {
    // Handle error
}
```

#### Executing Functions Within a Transaction

The `WithTx` method allows you to execute a function within the context of a transaction. It manages the lifecycle of the transaction automatically:

```go
err := txManager.WithTx(ctx, func(ctx context.Context, tx pgxpool.Tx) error {
    // Perform your transactional operations here
	
    _, err = tx.Query(ctx, 'query', 'params ...')
    if err != nil {
        return err
    }
	
    // ...
	
    return nil
})

if err != nil {
    // Handle error
}
```

For more complex scenarios involving nested transactions, use the `WithNestedTx` method:

```go
err := txManager.WithNestedTx(ctx, func(ctx context.Context, tx pgxpool.Tx) error {
    // Perform operations within a nested transaction here

    _, err = tx.Query(ctx, 'query', 'params ...')
    if err != nil {
        return err 
    }

    txManager.WithNestedTx(ctx, func(ctx context.Context, tx pgxpool.Tx) error {
        _, err = tx.Query(ctx, 'query2', 'params ...')
        if err != nil {
            return err
        }

        return nil
    })

    // ...
	
    return nil
})
if err != nil {
    // Handle error
}
```

## Migrate

`Migrate` is a Go package for managing database migrations, specifically designed for PostgreSQL using the `pgx` library.

### Features

- Easy setup and management of database migrations.
- Support for PostgreSQL databases using [golang-migrate](https://github.com/golang-migrate/migrate).
- Flexible configuration options through Viper or direct configuration.

### Using

### Creating a New Migrator

#### Using Direct Configuration

```go
import (
    "github.com/i4erkasov/go-pgsql/migrate"
)

func main() {
    config := migrate.Config{
        Driver:          "postgres",
        DataSourceName:  "postgres://user:password@localhost:5432/dbname",
        PathMigration:   "migrations",
        Scheme:          "public",
        Table:           "migration",
        ConnMaxLifetime: 10 * time.Minute,
    }

    migrator, err := migrate.NewWithConfig(config)
    if err != nil {
        // Handle error
    }
}
```

#### Using Viper for Configuration

```yaml
pgsql:
  migrations:
      driver: "postgres" # Default: "postgres"
      table: "migrations" # Default: "migrations"
      path: "migrations" # Default: "migrations"
      scheme: "public" # Default: "public"
      host: "127.0.0.1" # Required if not `dsn`
      port: "5432" # Required if not `dsn`
      user: "user" # Required if not `dsn`
      password: "password" # Required if not `dsn` 
      dbname: "db" # Required if not `dsn`
      conn_max_lifetime: '10m' # Default: 10 minutes
      options: # Optional
        sslmode: "disable" # Optional
        encoding: "UTF8" # Optional
```

#### or

```yaml
pgsql:
  migrations:
      driver: "postgres" # Default: "postgres"
      table: "migrations" # Default: "migrations"
      path: "migrations" # Default: "migrations"
      scheme: "public" # Default: "public"
      dsn: "postgres://user:password@127.0.0.1:5432/db?sslmode=disable&client_encoding=UTF8"
      conn_max_lifetime: '10m' # Default: 10 minutes
```

```go
import (
    "github.com/spf13/viper"
    "github.com/i4erkasov/go-pgsql/migrate"
)

func main() {
    viper.SetConfigName("config") // Name of config file (without extension)
    viper.AddConfigPath(".")      // Path to look for the config file in
    viper.AutomaticEnv()          // Automatically override values from environment variables

    migrator, err := migrate.NewWithViper(viper.GetViper())
    if err != nil {
        // Handle error
    }
}
```

#### Configuration Parameters

- `Driver`: Database driver (default: `"postgres"`).
- `DataSourceName`: Data source name or DSN (required).
- `PathMigration`: Path to migration files (default: `"migrations"`).
- `Scheme`: Database schema (default: `"public"`).
- `Table`: Migration table name (default: `"migration"`).
- `ConnMaxLifetime`: Maximum lifetime of database connections (default: `10 * time.Minute`).

#### Applying and Reverting Migrations

```go
err := migrator.Up(context.Background())
if err != nil {
    // Handle error
}
    
// To revert migrations
err = migrator.Down(context.Background(), 1) // Revert 1 step
if err != nil {
    // Handle error
}
```

#### An example of using migration with the [cobra](https://github.com/spf13/cobra) package

```go
import (
    "errors"
    "github.com/i4erkasov/go-pgsql/migrate"
    "github.com/spf13/cobra"
    "golang.org/x/exp/slices"
)

const (
    MigrationCommand      = "sql-migrate"
    VersionMigrateService = "0.0.1"
    MigrationUpArg        = "up"
    MigrationDownArg      = "down"
)

var steps int // variable to store the number of steps

var sqlMigrate = &cobra.Command{
    Use:        MigrationCommand,
    Short:      "Run database migration",
    Version:    VersionMigrateService,
    Args:       cobra.MaximumNArgs(1),
    ArgAliases: []string{MigrationUpArg, MigrationDownArg},
    RunE: func(cmd *cobra.Command, args []string) (err error) {
        var migration migrate.SQLMigrator
        if migration, err = migrate.NewWithViper(cfg.Sub("app")); err != nil {
            return err
        }

        switch true {
        case slices.Contains(args, MigrationUpArg):
            return migration.Up(cmd.Context())
        case slices.Contains(args, MigrationDownArg):
            return migration.Down(cmd.Context(), steps)
        default:
            return errors.New("invalid argument. please specify 'up' or 'down' for migration")
        }
    },
}

func init() {
    sqlMigrate.Flags().IntVarP(&steps, "steps", "s", 1, "Number of steps to migrate down")

    cmd.AddCommand(sqlMigrate)
}
```

## Contributing

Contributions to the `go-pgsql` package are welcome. Please submit pull requests to [GitHub repository](https://github.com/i4erkasov/go-pgsql).

Licensed under the [MIT License](LICENSE).
