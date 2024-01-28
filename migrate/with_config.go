package migrate

import (
	"errors"
	"fmt"
)

// NewWithConfig creates a new instance of Migrate using the provided Config.
func NewWithConfig(config Config) (*Migrate, error) {
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	SetDefaultConfigValues(&config)

	return &Migrate{
		config: config,
	}, nil
}

// SetDefaultConfigValues sets default values for Config fields if they are not set.
func SetDefaultConfigValues(config *Config) {
	if config.PathMigration == "" {
		config.PathMigration = defaultPath
	}
	if config.Scheme == "" {
		config.Scheme = defaultScheme
	}
	if config.Table == "" {
		config.Table = defaultTable
	}
	if config.ConnMaxLifetime == 0 {
		config.ConnMaxLifetime = defaultConnMaxLifetime
	}
	if config.Driver == "" {
		config.Driver = defaultDriver
	}
}

// validateConfig checks if the provided Config has all required fields.
func validateConfig(config Config) error {
	// Check if DataSourceName is set
	if config.DataSourceName == "" {
		return errors.New("data_source_name is required in the config")
	}

	// Check if PathMigration is set
	if config.PathMigration == "" {
		return errors.New("path_migration is required in the config")
	}

	// Check if Scheme is set
	if config.Scheme == "" {
		return errors.New("scheme is required in the config")
	}

	// Check if Table is set
	if config.Table == "" {
		return errors.New("table is required in the config")
	}

	// Check if ConnMaxLifetime is set
	if config.ConnMaxLifetime <= 0 {
		return errors.New("conn_max_lifetime must be greater than zero")
	}

	return nil
}
