package logger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNew_ReturnsBuildError(t *testing.T) {
	previous := buildZapLogger
	buildZapLogger = func(cfg zap.Config) (*zap.Logger, error) {
		return nil, errors.New("build failed")
	}
	t.Cleanup(func() {
		buildZapLogger = previous
	})

	l, err := New(Config{Level: "info", Format: "json"})

	require.Error(t, err)
	assert.Nil(t, l)
	assert.Contains(t, err.Error(), "build logger")
}
