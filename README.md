# go-pgsql

`go-pgsql` is a comprehensive Go package that facilitates managing PostgreSQL connection pools, database migrations, and transaction management. It leverages the robust `pgx` library for pool management and `golang-migrate` for database migrations, providing a full suite of tools for working with PostgreSQL in Go applications.

## Features

- **Connection Pool Management (`pgxpool`)**: Simplifies the creation and management of connection pools using `pgxpool`.
- **Database Migrations**: Allows for smooth database schema migrations using the `golang-migrate` package.
- **Transaction Management**: Offers a convenient transaction manager that supports nested transactions and automatic error handling.

## External Packages

This package uses the following external libraries:

- Connection Pooling: [github.com/jackc/pgx/v4/pgxpool](https://github.com/jackc/pgx)
- Database Migrations: [github.com/golang-migrate/migrate/v4](https://github.com/golang-migrate/migrate)
- Configuration: [github.com/spf13/viper](https://github.com/spf13/viper)

## Installation

```bash
go get github.com/i4erkasov/go-pgsql
```

## Pool Registry

A package for managing PostgreSQL connection pools using the [pgxpool](https://github.com/jackc/pgx) library. It offers a convenient way to set up and manage database connections for different environments.

### Features

- Easy creation and management of PostgreSQL connection pools.
- Supports various configurations for different pools.
- Seamless integration with configuration systems like Viper.
- Flexible configuration options through functional options pattern.

### Usage the `pgxpool`
#### Creating a New Pool Registry:

#### Using Viper for Configuration

```yaml
pgsql:
  pgpool:
    default:
      nodes: [
        "postgres://user:password@127.0.0.1:5432/db?sslmode=disable&client_encoding=UTF8"
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

    pool, err := registry.Pools() // is DEFAULT pool getter
    pool, err := registry.GetPoolName("pool1") // is pool getter by name.
}
```
### or

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


## Transaction Management

The package includes a sophisticated transaction manager that allows for simple and complex transactional operations, including support for nested transactions.

### Features

- Easy to use transaction manager with support for context-based transactions.
- Support for nested transactions, allowing complex transactional workflows.
- Automatic rollback on errors and panics to maintain data integrity.
- Simple interface to execute transactional functions.

### Using the `TxManager`

#### Starting Transactions

To begin a transaction, simply call the `Begin` method on an instance of `TxManager`:

```go
import (
    "github.com/i4erkasov/go-pgsql/tx"
)

// ...

txManager, err := tx.NewTxManager(registry)
if err != nil {
// Handle error
}

ctx, err := txManager.Begin(context.Background())
if err != nil {
// Handle error
}
```

#### Committing Transactions
Once your operations are complete, `Commit` the transaction:

```go
err = txManager.Commit(ctx)
if err != nil {
    // Handle error
}
```

#### Rolling Back Transactions
If an error occurs, or you need to abort the transaction, call `Rollback`:

```go
err = txManager.Rollback(ctx)
if err != nil {
    // Handle error
}
```

#### Executing Functions Within a Transaction

The `WithTx` method allows you to execute a function within the context of a transaction. It manages the lifecycle of the transaction automatically:

```go
err := txManager.WithTx(ctx, func(ctx context.Context, tx pgxpool.Tx) error {
    // Perform your transactional operations here
    return nil
})
if err != nil {
    // Handle error
}
```

Nested Transactions
For more complex scenarios involving nested transactions, use the `WithNestedTx` method:

```go
err := txManager.WithNestedTx(ctx, func(ctx context.Context, tx pgxpool.Tx) error {
    // Perform operations within a nested transaction here
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

### Usage the `Migrator`

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

## Contributing

Contributions to the `go-pgsql` package are welcome. Please submit pull requests to [GitHub repository](https://github.com/i4erkasov/go-pgsql).
