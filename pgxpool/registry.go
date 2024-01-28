package pgxpool

import (
	"errors"
	"sync"
	"time"
)

const (
	// DEFAULT is default pool name.
	DEFAULT                  = "default"
	defaultMaxConns          = int32(4)
	defaultMinConns          = int32(0)
	defaultMaxConnLifetime   = time.Hour
	defaultMaxConnIdleTime   = time.Minute * 30
	defaultHealthCheckPeriod = time.Minute
)

type (
	// Configs is registry configurations.
	Configs map[string]Config

	// Config single config item
	Config struct {
		Nodes                []string      `mapstructure:"nodes" json:"nodes"`
		MaxConns             int32         `mapstructure:"max_conns" json:"max_conns"`
		MinConns             int32         `mapstructure:"min_conns" json:"min_conns"`
		MaxConnLifetime      time.Duration `mapstructure:"max_conn_lifetime" json:"max_conn_lifetime"`
		MaxConnIdleTime      time.Duration `mapstructure:"max_conn_idle_time" json:"max_conn_idle_time"`
		HealthCheckPeriod    time.Duration `mapstructure:"health_check_period" json:"health_check_period"`
		LazyConnect          bool          `mapstructure:"lazy_conn" json:"lazy_conn"`
		PreferSimpleProtocol bool          `mapstructure:"prefer_simple_protocol" json:"prefer_simple_protocol"`
	}

	// Registry is database pool registry.
	Registry struct {
		sync.Mutex
		pools map[string]*Pools
		conf  Configs
	}
)

var (
	// ErrUnknownPool is error triggered when pool with provided name not founded.
	ErrUnknownPool = errors.New("unknown pool")
)

// NewRegistry is registry constructor.
func NewRegistry(configs Configs) (*Registry, error) {
	pools := make(map[string]*Pools, len(configs))

	for name, config := range configs {
		p, err := Open(config)
		if err != nil {
			return nil, err
		}

		pools[name] = p
	}

	return &Registry{
		pools: pools,
		conf:  configs,
	}, nil
}

func GetDefaultConfig() Config {
	return Config{
		MaxConns:          defaultMaxConns,
		MinConns:          defaultMinConns,
		MaxConnLifetime:   defaultMaxConnLifetime,
		MaxConnIdleTime:   defaultMaxConnIdleTime,
		HealthCheckPeriod: defaultHealthCheckPeriod,
	}
}

// SetDefaultValues sets default values for Config if they are not specified.
func SetDefaultValues(cfg *Config) {
	if cfg.MaxConns == 0 {
		cfg.MaxConns = defaultMaxConns
	}
	if cfg.MinConns == 0 {
		cfg.MinConns = defaultMinConns
	}
	if cfg.MaxConnLifetime == 0 {
		cfg.MaxConnLifetime = defaultMaxConnLifetime
	}
	if cfg.MaxConnIdleTime == 0 {
		cfg.MaxConnIdleTime = defaultMaxConnIdleTime
	}
	if cfg.HealthCheckPeriod == 0 {
		cfg.HealthCheckPeriod = defaultHealthCheckPeriod
	}
}

// Close is method for close pools connections.
func (r *Registry) Close() error {
	r.Lock()
	defer r.Unlock()

	for key, pool := range r.pools {
		pool.Close()
		delete(r.pools, key)
	}

	return nil
}

// Pools is default pool getter.
func (r *Registry) Pools() (*Pools, error) {
	return r.GetPoolName(DEFAULT)
}

// GetPoolName PoolWithName is pool getter by name.
func (r *Registry) GetPoolName(name string) (*Pools, error) {
	r.Lock()
	defer r.Unlock()

	if pool, ok := r.pools[name]; ok {
		return pool, nil
	}

	return nil, ErrUnknownPool
}

// mergeConfigs merges two Config objects. Non-zero values in cfg2 override values in cfg1.
func merge(old Config, new Config) Config {
	if new.MaxConns != 0 {
		old.MaxConns = new.MaxConns
	}
	if new.MinConns != 0 {
		old.MinConns = new.MinConns
	}
	if new.MaxConnLifetime != 0 {
		old.MaxConnLifetime = new.MaxConnLifetime
	}
	if new.MaxConnIdleTime != 0 {
		old.MaxConnIdleTime = new.MaxConnIdleTime
	}
	if new.HealthCheckPeriod != 0 {
		old.HealthCheckPeriod = new.HealthCheckPeriod
	}
	if len(new.Nodes) > 0 {
		old.Nodes = new.Nodes
	}
	old.LazyConnect = new.LazyConnect
	old.PreferSimpleProtocol = new.PreferSimpleProtocol

	return old
}
