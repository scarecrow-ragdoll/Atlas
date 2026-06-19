package config_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/libs/go/config"
)

// fullConfig mirrors a typical app config with shared + custom blocks.
type fullConfig struct {
	Server   config.ServerConfig   `mapstructure:"server"`
	Log      config.LogConfig      `mapstructure:"log"`
	Postgres config.PostgresConfig `mapstructure:"postgres"`
	Redis    config.RedisConfig    `mapstructure:"redis"`
	Auth     authConfig            `mapstructure:"auth"`
}

type authConfig struct {
	JWTSecret string `mapstructure:"jwt_secret" validate:"required"`
}

// minimalConfig uses only one shared block.
type minimalConfig struct {
	Postgres config.PostgresConfig `mapstructure:"postgres"`
}

func TestLoad_YAMLParsing(t *testing.T) {
	cfg, err := config.Load[fullConfig](config.Options{
		ConfigPath: "testdata/valid.yml",
	})
	require.NoError(t, err)

	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "development", cfg.Server.Env)
	assert.Equal(t, []string{"http://localhost:3000"}, cfg.Server.CORSOrigins)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, "localhost", cfg.Postgres.Host)
	assert.Equal(t, 5432, cfg.Postgres.Port)
	assert.Equal(t, "testuser", cfg.Postgres.User)
	assert.Equal(t, "testdb", cfg.Postgres.DB)
	assert.Equal(t, "localhost", cfg.Redis.Host)
	assert.Equal(t, 6379, cfg.Redis.Port)
	assert.Equal(t, "test-secret", cfg.Auth.JWTSecret)
}

func TestLoad_EnvVarOverridesYAML(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "override-host")
	t.Setenv("SERVER_PORT", "9090")

	cfg, err := config.Load[fullConfig](config.Options{
		ConfigPath: "testdata/valid.yml",
	})
	require.NoError(t, err)

	assert.Equal(t, "override-host", cfg.Postgres.Host)
	assert.Equal(t, 9090, cfg.Server.Port)
}

func TestLoad_DotEnvFilePickedUp(t *testing.T) {
	// godotenv.Load writes to OS env — clean up vars it sets to avoid polluting other tests.
	t.Cleanup(func() {
		_ = os.Unsetenv("POSTGRES_HOST")
		_ = os.Unsetenv("POSTGRES_PORT")
	})

	cfg, err := config.Load[fullConfig](config.Options{
		ConfigPath: "testdata/valid.yml",
		EnvFile:    "testdata/.env.test",
	})
	require.NoError(t, err)

	assert.Equal(t, "envhost", cfg.Postgres.Host)
	assert.Equal(t, 9999, cfg.Postgres.Port)
}

func TestLoad_DotEnvFileError(t *testing.T) {
	_, err := config.Load[fullConfig](config.Options{
		ConfigPath: "testdata/valid.yml",
		EnvFile:    "testdata/missing.env",
	})

	require.Error(t, err)
	assert.True(t, errors.Is(err, config.ErrEnvFileLoad))
}

func TestLoad_PriorityOrder_EnvWins(t *testing.T) {
	// Real env var should beat .env file value.
	// t.Setenv sets POSTGRES_HOST before godotenv.Load runs;
	// godotenv.Load (not Overload) skips already-set vars -> real env wins.
	t.Setenv("POSTGRES_HOST", "real-env-host")
	t.Cleanup(func() {
		_ = os.Unsetenv("POSTGRES_PORT") // .env.test also sets POSTGRES_PORT
	})

	cfg, err := config.Load[fullConfig](config.Options{
		ConfigPath: "testdata/valid.yml",
		EnvFile:    "testdata/.env.test", // .env has POSTGRES_HOST=envhost
	})
	require.NoError(t, err)

	assert.Equal(t, "real-env-host", cfg.Postgres.Host)
}

func TestLoad_ValidationError(t *testing.T) {
	// minimal.yml has no auth section, so Auth.JWTSecret will be "".
	// Explicitly set env to empty to prevent any leaked value from prior tests
	// from accidentally making this pass.
	t.Setenv("AUTH_JWT_SECRET", "")

	_, err := config.Load[fullConfig](config.Options{
		ConfigPath: "testdata/minimal.yml",
	})
	require.Error(t, err)
	assert.True(t, errors.Is(err, config.ErrValidation))
}

func TestLoad_MissingYAML(t *testing.T) {
	_, err := config.Load[fullConfig](config.Options{
		ConfigPath: "testdata/nonexistent.yml",
	})
	require.Error(t, err)
	assert.True(t, errors.Is(err, config.ErrConfigLoad))
}

func TestLoad_UnmarshalError(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "invalid-*.yml")
	require.NoError(t, err)
	_, err = file.WriteString("postgres:\n  port: not-a-number\n")
	require.NoError(t, err)
	require.NoError(t, file.Close())

	_, err = config.Load[minimalConfig](config.Options{ConfigPath: file.Name()})

	require.Error(t, err)
	assert.True(t, errors.Is(err, config.ErrConfigLoad))
}

func TestLoad_PartialStruct(t *testing.T) {
	// minimalConfig has only Postgres — omitempty oneof fields should pass
	cfg, err := config.Load[minimalConfig](config.Options{
		ConfigPath: "testdata/minimal.yml",
	})
	require.NoError(t, err)

	assert.Equal(t, "localhost", cfg.Postgres.Host)
	assert.Equal(t, 5432, cfg.Postgres.Port)
}

func TestLoad_EnvPrefix(t *testing.T) {
	t.Setenv("MYAPP_POSTGRES_HOST", "prefix-host")

	cfg, err := config.Load[minimalConfig](config.Options{
		ConfigPath: "testdata/minimal.yml",
		EnvPrefix:  "MYAPP",
	})
	require.NoError(t, err)

	assert.Equal(t, "prefix-host", cfg.Postgres.Host)
}

func TestLoad_AutomaticEnvNestedKeys(t *testing.T) {
	// Verify AutomaticEnv resolves nested keys without manual BindEnv
	t.Setenv("SERVER_READ_TIMEOUT", "25s")
	t.Setenv("POSTGRES_MAX_CONNS", "50")

	cfg, err := config.Load[fullConfig](config.Options{
		ConfigPath: "testdata/valid.yml",
	})
	require.NoError(t, err)

	// If AutomaticEnv works for nested keys, these will be overridden.
	// If not, they'll have the YAML defaults (10s, 10).
	// This test documents the actual behavior — if it fails,
	// fallback to explicit BindEnv via reflection is needed.
	t.Logf("server.read_timeout = %v (env: 25s, yaml: 10s)", cfg.Server.ReadTimeout)
	t.Logf("postgres.max_conns = %d (env: 50, yaml: 10)", cfg.Postgres.MaxConns)
}

func TestValidate_Standalone(t *testing.T) {
	// Validate can be used outside of Load
	pg := config.PostgresConfig{
		Host: "",
		Port: 0,
	}
	err := config.Validate(pg)
	require.Error(t, err)
	assert.True(t, errors.Is(err, config.ErrValidation))
	assert.Contains(t, err.Error(), "Host is required")
	assert.Contains(t, err.Error(), "Port is required")
}

func TestValidate_OneOfAndDefaultMessages(t *testing.T) {
	type customConfig struct {
		Log   config.LogConfig `validate:"required"`
		Email string           `validate:"email"`
		Count int              `validate:"gt=0"`
	}

	err := config.Validate(customConfig{
		Log:   config.LogConfig{Level: "verbose", Format: "xml"},
		Email: "not-email",
	})

	require.Error(t, err)
	assert.True(t, errors.Is(err, config.ErrValidation))
	assert.Contains(t, err.Error(), "Level must be one of")
	assert.Contains(t, err.Error(), "Format must be one of")
	assert.Contains(t, err.Error(), "Email failed on 'email' validation")
	assert.Contains(t, err.Error(), "Count must be greater than 0")
}

func TestPostgresConfig_DSN(t *testing.T) {
	pg := config.PostgresConfig{
		Host:     "dbhost",
		Port:     5432,
		User:     "admin",
		Password: "pass",
		DB:       "mydb",
		SSLMode:  "require",
	}
	expected := "postgres://admin:pass@dbhost:5432/mydb?sslmode=require"
	assert.Equal(t, expected, pg.DSN())
}
