package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate_WrapsNonValidationError(t *testing.T) {
	previous := validateStruct
	validateStruct = func(s any) error {
		return errors.New("validator unavailable")
	}
	t.Cleanup(func() {
		validateStruct = previous
	})

	err := Validate(struct{}{})

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrValidation)
	assert.Contains(t, err.Error(), "validator unavailable")
}
