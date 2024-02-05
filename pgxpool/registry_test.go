package pgxpool

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RegistryTestSuite struct {
	suite.Suite
}

// TestSetDefaultValues checks if default values are correctly set.
func (t *RegistryTestSuite) TestSetDefaultValues() {
	t.T().Parallel()

	cfg := Config{}                  // Create an empty config.
	SetDefaultValues(&cfg)           // Set default values.
	defaultCfg := GetDefaultConfig() // Get the default config for comparison.

	t.Equal(defaultCfg.MaxConns, cfg.MaxConns, "MaxConns should match the default value")
	t.Equal(defaultCfg.MinConns, cfg.MinConns, "MinConns should match the default value")
	t.Equal(defaultCfg.MaxConnLifetime, cfg.MaxConnLifetime, "MaxConnLifetime should match the default value")
	t.Equal(defaultCfg.MaxConnIdleTime, cfg.MaxConnIdleTime, "MaxConnIdleTime should match the default value")
	t.Equal(defaultCfg.HealthCheckPeriod, cfg.HealthCheckPeriod, "HealthCheckPeriod should match the default value")
}

// TestGetPoolNameNotFound checks that an error is returned for a non-existent pool name.
func (t *RegistryTestSuite) TestGetPoolNameNotFound() {
	t.T().Parallel()

	registry, err := NewRegistry(Configs{})
	t.NoError(err, "Creating registry should not produce an error")

	_, err = registry.GetPoolName("nonexistent")
	t.Error(err, "Should return an error for a non-existent pool name")
	t.Equal(ErrUnknownPool, err, "Error should be ErrUnknownPool for a non-existent pool name")
}

// RegistrySuite runs the test suite.
func TestRegistrySuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(RegistryTestSuite))
}
