package migrate

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var required = []string{"user", "password", "host", "port", "dbname"}

// NewWithViper creates a new instance of a migration.
func NewWithViper(cfg *viper.Viper) (SQLMigrator, error) {
	cfg = cfg.Sub(cfgParamName)
	if cfg == nil {
		return nil, ErrMigrationConfigMissing
	}

	if !cfg.IsSet("dsn") {
		if err := validate(cfg); err != nil {
			return nil, err
		}
	}

	dsn := cfg.GetString("dsn")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			cfg.GetString("user"),
			cfg.GetString("password"),
			cfg.GetString("host"),
			cfg.GetString("port"),
			cfg.GetString("dbname"),
		)
		// Add additional parameters
		var params []string
		if cfg.IsSet("options.sslmode") {
			params = append(params, fmt.Sprintf("sslmode=%s", cfg.GetString("options.sslmode")))
		}

		if cfg.IsSet("options.encoding") {
			params = append(params, fmt.Sprintf("client_encoding=%s", cfg.GetString("options.encoding")))
		}

		if len(params) > 0 {
			dsn += "?" + strings.Join(params, "&")
		}
	}

	cfg.SetDefault("path", defaultPath)
	cfg.SetDefault("driver", defaultDriver)
	cfg.SetDefault("schema", defaultScheme)
	cfg.SetDefault("table", defaultTable)
	cfg.SetDefault("conn_max_lifetime", defaultConnMaxLifetime)

	return &Migrate{
		config: Config{
			DataSourceName:  os.ExpandEnv(dsn),
			Driver:          cfg.GetString("driver"),
			PathMigration:   cfg.GetString("path"),
			Scheme:          cfg.GetString("schema"),
			Table:           cfg.GetString("table"),
			ConnMaxLifetime: cfg.GetDuration("conn_max_lifetime"),
		},
	}, nil
}

// validate checks if all required fields are present in the migration configuration.
func validate(cfg *viper.Viper) error {
	for _, field := range required {
		if !cfg.IsSet(field) {
			return fmt.Errorf("required field '%s' is missing in migration configuration", field)
		}
	}
	return nil
}
