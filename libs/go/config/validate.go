package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())
var validateStruct = validate.Struct

// Validate runs struct-tag validation on v and returns a human-readable error.
func Validate(v any) error {
	err := validateStruct(v)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}

	msgs := make([]string, 0, len(validationErrors))
	for _, fe := range validationErrors {
		msgs = append(msgs, formatFieldError(fe))
	}

	return fmt.Errorf("%w: %s", ErrValidation, strings.Join(msgs, "; "))
}

func formatFieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s failed on '%s' validation", fe.Field(), fe.Tag())
	}
}
