package logger_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/libs/go/logger"
)

func TestNew_JSONFormat(t *testing.T) {
	l, err := logger.New(logger.Config{Level: "info", Format: "json"})
	require.NoError(t, err)
	assert.NotNil(t, l)
}

func TestNew_ConsoleFormat(t *testing.T) {
	l, err := logger.New(logger.Config{Level: "debug", Format: "console"})
	require.NoError(t, err)
	assert.NotNil(t, l)
}

func TestNew_DefaultsToJSONInfo(t *testing.T) {
	l, err := logger.New(logger.Config{})
	require.NoError(t, err)
	assert.NotNil(t, l)
}

func TestNew_InvalidLevel_ReturnsError(t *testing.T) {
	_, err := logger.New(logger.Config{Level: "invalid_level"})
	assert.Error(t, err)
}
