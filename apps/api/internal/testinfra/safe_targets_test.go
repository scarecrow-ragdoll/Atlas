package testinfra

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeTestingT struct {
	fatalMessage string
}

func (f *fakeTestingT) Helper() {}

func (f *fakeTestingT) Fatalf(format string, args ...any) {
	if f.fatalMessage == "" {
		f.fatalMessage = fmt.Sprintf(format, args...)
	}
}

func TestCoverageGateEnabledReadsEnv(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		t.Setenv("COVERAGE_GATE", "1")

		assert.True(t, CoverageGateEnabled())
	})

	t.Run("false when unset", func(t *testing.T) {
		t.Setenv("COVERAGE_GATE", "")

		assert.False(t, CoverageGateEnabled())
	})

	t.Run("false for other value", func(t *testing.T) {
		t.Setenv("COVERAGE_GATE", "true")

		assert.False(t, CoverageGateEnabled())
	})
}

func TestValidateSafePostgresDSNAcceptsDefaultTestTarget(t *testing.T) {
	t.Setenv("TEST_POSTGRES_DB", "monorepo_test")
	t.Setenv("TEST_POSTGRES_PORT", "17501")

	err := ValidateSafePostgresDSN("postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable")

	require.NoError(t, err)
}

func TestValidateSafePostgresDSNRejectsUnsafeTargets(t *testing.T) {
	t.Setenv("TEST_POSTGRES_DB", "monorepo_test")
	t.Setenv("TEST_POSTGRES_PORT", "17501")
	// #nosec G101 -- local Docker test fixture credentials.
	cases := map[string]string{
		"empty":       "",
		"malformed":   "://bad",
		"wrongScheme": "mysql://app:secret@localhost:17501/monorepo_test",
		"devDB":       "postgres://app:secret@localhost:17501/monorepo_dev?sslmode=disable",
		"wrongDB":     "postgres://app:secret@localhost:17501/monorepo_other?sslmode=disable",
		"devPort":     "postgres://app:secret@localhost:7501/monorepo_test?sslmode=disable",
		"wrongPort":   "postgres://app:secret@localhost:15432/monorepo_test?sslmode=disable",
		"missingPort": "postgres://app:secret@localhost/monorepo_test?sslmode=disable",
	}

	for name, dsn := range cases {
		t.Run(name, func(t *testing.T) {
			err := ValidateSafePostgresDSN(dsn)

			require.Error(t, err)
		})
	}
}

func TestValidateSafePostgresDSNRejectsDevDatabaseEvenWhenConfiguredAsExpected(t *testing.T) {
	t.Setenv("TEST_POSTGRES_DB", "monorepo_dev")
	t.Setenv("TEST_POSTGRES_PORT", "17501")

	err := ValidateSafePostgresDSN("postgres://app:secret@localhost:17501/monorepo_dev?sslmode=disable")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "development postgres database")
}

func TestRequireSafePostgresDSNFailsUnsafeTarget(t *testing.T) {
	fakeT := &fakeTestingT{}

	RequireSafePostgresDSN(fakeT, "postgres://app:secret@localhost:7501/monorepo_dev?sslmode=disable")

	require.NotEmpty(t, fakeT.fatalMessage)
	assert.Contains(t, fakeT.fatalMessage, "unsafe postgres test DSN")
}

func TestPostgresDSNDefaultsToTestDatabase(t *testing.T) {
	t.Setenv("API_TEST_DATABASE_DSN", "")
	t.Setenv("TEST_POSTGRES_HOST", "127.0.0.1")
	t.Setenv("TEST_POSTGRES_PORT", "17501")
	t.Setenv("TEST_POSTGRES_USER", "app")
	t.Setenv("TEST_POSTGRES_PASSWORD", "secret")
	t.Setenv("TEST_POSTGRES_DB", "monorepo_test")

	dsn := PostgresDSN()

	assert.Equal(t, "postgres://app:secret@127.0.0.1:17501/monorepo_test?sslmode=disable", dsn)
}

func TestPostgresDSNUsesEnvOverride(t *testing.T) {
	t.Setenv("API_TEST_DATABASE_DSN", "postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable")

	dsn := PostgresDSN()

	assert.Equal(t, "postgres://app:secret@localhost:17501/monorepo_test?sslmode=disable", dsn)
}

func TestPostgresConfigUsesTestDefaults(t *testing.T) {
	t.Setenv("TEST_POSTGRES_HOST", "localhost")
	t.Setenv("TEST_POSTGRES_PORT", "17501")
	t.Setenv("TEST_POSTGRES_USER", "app")
	t.Setenv("TEST_POSTGRES_PASSWORD", "secret")
	t.Setenv("TEST_POSTGRES_DB", "monorepo_test")

	cfg := PostgresConfig(t)

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 17501, cfg.Port)
	assert.Equal(t, "app", cfg.User)
	assert.Equal(t, "secret", cfg.Password)
	assert.Equal(t, "monorepo_test", cfg.DB)
	assert.Equal(t, "disable", cfg.SSLMode)
}

func TestPostgresConfigRejectsNonPositivePort(t *testing.T) {
	t.Setenv("TEST_POSTGRES_PORT", "0")

	fakeT := &fakeTestingT{}
	_ = PostgresConfig(fakeT)

	require.NotEmpty(t, fakeT.fatalMessage)
	assert.Contains(t, fakeT.fatalMessage, "greater than zero")
}

func TestRedisConfigUsesTestDefaults(t *testing.T) {
	t.Setenv("TEST_REDIS_HOST", "localhost")
	t.Setenv("TEST_REDIS_PORT", "17502")
	t.Setenv("TEST_REDIS_PASSWORD", "")
	t.Setenv("TEST_REDIS_DB", "0")

	cfg := RedisConfig(t)

	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 17502, cfg.Port)
	assert.Equal(t, "", cfg.Password)
	assert.Equal(t, 0, cfg.DB)
}

func TestRedisConfigRejectsDevelopmentPort(t *testing.T) {
	t.Setenv("TEST_REDIS_PORT", "7502")

	fakeT := &fakeTestingT{}
	_ = RedisConfig(fakeT)

	require.NotEmpty(t, fakeT.fatalMessage)
	assert.True(t, strings.Contains(fakeT.fatalMessage, "development redis port"), fakeT.fatalMessage)
}

func TestRedisConfigRejectsInvalidDatabase(t *testing.T) {
	t.Setenv("TEST_REDIS_PORT", "17502")
	t.Setenv("TEST_REDIS_DB", "not-a-number")

	fakeT := &fakeTestingT{}
	_ = RedisConfig(fakeT)

	require.NotEmpty(t, fakeT.fatalMessage)
	assert.Contains(t, fakeT.fatalMessage, "must be an integer")
}
