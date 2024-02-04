package pgxpool

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)

// Pool is an alias to pgxpool.Pool
type Pool = pgxpool.Pool

const cfgParamName = "pgsql.pgpool"

func NewWithViper(cfg *viper.Viper) (*Registry, error) {
	var (
		keys   = cfg.Sub(cfgParamName).AllKeys()
		config = make(Configs, len(keys))
	)

	for _, key := range keys {
		var name = strings.Split(key, ".")[0]
		if _, ok := config[name]; ok {
			continue
		}

		k := fmt.Sprintf("%s.%s", cfgParamName, name)

		if !cfg.IsSet(k + ".nodes") {
			return nil, fmt.Errorf("config key \"%s\" is required", k+"nodes")
		}

		conf := GetDefaultConfig()
		if err := cfg.Sub(k).Unmarshal(&conf); err != nil {
			return nil, err
		}

		config[name] = conf
	}

	return NewRegistry(config)
}

// NewWithConfigOptions creates a new Registry directly using the provided configuration options.
func NewWithConfigOptions(opts ...ConfigOption) (*Registry, error) {
	configs := make(map[string]Config)

	// Apply each provided option to the configs
	for _, opt := range opts {
		opt(configs)
	}

	return NewWithConfigs(configs)
}

// NewWithConfigs creates a new Registry using the provided configurations.
func NewWithConfigs(configs map[string]Config) (*Registry, error) {
	// Setting default values for each configuration if not specified
	for name, cfg := range configs {
		SetDefaultValues(&cfg)
		configs[name] = cfg
	}

	return NewRegistry(configs)
}

// ConfigOption defines the type for functional options for Registry configuration.
type ConfigOption func(map[string]Config)

// WithConfig is an option to add a configuration to the Registry.
func WithConfig(name string, cfg Config) ConfigOption {
	return func(configs map[string]Config) {
		if existingCfg, exists := configs[name]; exists {
			cfg = merge(existingCfg, cfg)
		} else {
			SetDefaultValues(&cfg)
		}
		configs[name] = cfg
	}
}
