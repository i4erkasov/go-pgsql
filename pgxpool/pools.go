package pgxpool

import (
	"context"
	"sync/atomic"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Pools wraps pgx.Pools for master/slave support
type Pools struct {
	pools []*Pool
	count uint64
}

// Open creates new pools for each node
func Open(config Config) (*Pools, error) {
	pools := make([]*Pool, 0, len(config.Nodes))

	for _, node := range config.Nodes {
		c, err := pgxpool.ParseConfig(node)
		if err != nil {
			return nil, err
		}

		c.MaxConns = config.MaxConns
		c.MinConns = config.MinConns
		c.MaxConnLifetime = config.MaxConnLifetime
		c.MaxConnIdleTime = config.MaxConnIdleTime
		c.HealthCheckPeriod = config.HealthCheckPeriod
		c.LazyConnect = config.LazyConnect

		c.ConnConfig.PreferSimpleProtocol = config.PreferSimpleProtocol
		if config.PreferSimpleProtocol {
			c.ConnConfig.RuntimeParams["standard_conforming_strings"] = "on"
		}

		pool, err := pgxpool.ConnectConfig(context.Background(), c)
		if err != nil {
			return nil, err
		}

		pools = append(pools, pool)
	}

	return &Pools{pools: pools}, nil
}

// Close closes all connections in the pool and rejects future Acquire calls
func (p *Pools) Close() {
	for _, pool := range p.pools {
		pool.Close()
	}
}

// Master returns master connections pool
func (p *Pools) Master() *Pool {
	return p.pools[0]
}

// Slave returns slave connections pool
func (p *Pools) Slave() *Pool {
	return p.pools[p.slave(len(p.pools))]
}

func (p *Pools) slave(n int) int {
	if n <= 1 {
		return 0
	}

	return int(1 + (atomic.AddUint64(&p.count, 1) % uint64(n-1)))
}
