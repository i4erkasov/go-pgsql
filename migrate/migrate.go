package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"time"
)

const (
	cfgParamName = "pgsql.migrations"
	// defaultDriver is default driver
	defaultDriver = "postgres"
	// defaultPath is default path for migrations
	defaultPath = "migrations"
	// defaultScheme is default scheme in DB
	defaultScheme = "public"
	// defaultTable is default table in DB
	defaultTable = "migration"
	// defaultConnMaxLifetime is default maximum time for connecting to the database
	defaultConnMaxLifetime = 10 * time.Minute
)

var ErrMigrationConfigMissing = errors.New("migration configuration is missing")

// Config is registry configuration item.
type Config struct {
	Driver          string
	DatabaseName    string
	DataSourceName  string
	PathMigration   string
	Scheme          string
	Table           string
	ConnMaxLifetime time.Duration
}

type Migrate struct {
	config Config
}

// SQLMigrator works with migrations
type SQLMigrator interface {
	Up(ctx context.Context) error
	Down(ctx context.Context, max int) error
}

// Up applies all available migrations.
func (m *Migrate) Up(ctx context.Context) error {
	mig, err := m.newMigrationInstance(ctx)
	if err != nil {
		return fmt.Errorf("failed to create a new migration instance: %w", err)
	}

	defer func() {
		if _, dbErr := mig.Close(); dbErr != nil {
			if err == nil {
				err = fmt.Errorf("failed to close migration: %w", dbErr)
			}
		}
	}()

	err = mig.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// Down reverts a specified number of migrations.
func (m *Migrate) Down(ctx context.Context, max int) error {
	mig, err := m.newMigrationInstance(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if _, dbErr := mig.Close(); dbErr != nil {
			if err == nil {
				err = fmt.Errorf("failed to close migration: %w", dbErr)
			}
		}
	}()

	err = mig.Steps(-max)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

// openConnection opens a new database connection.
func (m *Migrate) openConnection() (*sql.DB, error) {
	db, err := sql.Open(m.config.Driver, m.config.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)

	}
	db.SetConnMaxLifetime(m.config.ConnMaxLifetime)
	return db, nil
}

// closeConnection closes the given database connection.
func (m *Migrate) closeConnection(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close the database connection: %w", err)
	}
	return nil
}

// newMigrationInstance creates a new migration instance with a managed database connection.
func (m *Migrate) newMigrationInstance(ctx context.Context) (_ *migrate.Migrate, err error) {
	db, err := m.openConnection()
	if err != nil {
		return nil, err
	}

	// Ensure the database connection is closed when the function returns
	defer func() {
		if closeErr := m.closeConnection(db); closeErr != nil {
			if err == nil {
				err = fmt.Errorf("failed to close the database connection: %w", closeErr)
			}
		}
	}()

	// Close the connection when context is cancelled
	go func() {
		<-ctx.Done()
		_ = m.closeConnection(db) // Ignore the error
	}()

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to establish a connection with the database: %w", err)
	}

	if err = m.checkMigrationTable(db); err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: m.config.Table,
		SchemaName:      m.config.Scheme,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create migration driver: %w", err)
	}

	mig, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", m.config.PathMigration), m.config.DatabaseName, driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	return mig, nil
}

// checkMigrationTable checks for the existence of the migration table and creates it if necessary.
func (m *Migrate) checkMigrationTable(db *sql.DB) error {
	var tableName string
	// Query to check if the migration table exists in the specified schema and table name
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = $1 AND table_name = $2"
	err := db.QueryRow(query, m.config.Scheme, m.config.Table).Scan(&tableName)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		// Table not found, create it
		query = "CREATE TABLE $1.$2 (version bigint not null primary key, dirty boolean not null)"
		_, err = db.Exec(query, m.config.Scheme, m.config.Table)
		if err != nil {
			return fmt.Errorf("failed to create migration table: %w", err)
		}
		return nil
	case err != nil:
		// An error occurred during the query
		return fmt.Errorf("failed to query migration table: %w", err)
	default:
		// Table exists
		return nil
	}
}
