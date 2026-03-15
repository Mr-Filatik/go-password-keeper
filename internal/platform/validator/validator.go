// Package validator provides functionality for validators.
package validator

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// IValidator provides a generic interface for validators.
type IValidator interface {
	// Validate checks that the model meets the constraints.
	//
	// Returns an error and string messages for each validation issue.
	Validate(entity any) ([]string, error)
}

// Validator represents a validator that uses the github.com/go-playground/validator/v10 package internally.
type Validator struct {
	validator *validator.Validate
}

// New returns a new validator instance.
func New() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

// Validate checks that the model meets the constraints.
//
// Returns an error and string messages for each validation issue.
//
// Implements the IValidator interface.
func (v *Validator) Validate(entity any) ([]string, error) {
	err := v.validator.Struct(entity)
	if err != nil {
		var problems []string

		var validationErrors validator.ValidationErrors

		if errors.As(err, &validationErrors) {
			for _, fieldError := range validationErrors {
				problems = append(problems, getValidationMessage(fieldError))
			}
		}

		return problems, fmt.Errorf("validate failed: %w", err)
	}

	return []string{}, nil
}

//nolint:funlen,cyclop
func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("Field '%s' is required", err.Field())

	case "min":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("Field '%s' must be at least %s characters long",
				err.Field(), err.Param())
		}

		return fmt.Sprintf("Field '%s' must be >= %s", err.Field(), err.Param())

	case "max":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("Field '%s' must be at most %s characters long",
				err.Field(), err.Param())
		}

		return fmt.Sprintf("Field '%s' must be <= %s", err.Field(), err.Param())

	case "gte":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("Field '%s' must be at least %s characters long",
				err.Field(), err.Param())
		}

		return fmt.Sprintf("Field '%s' must be >= %s", err.Field(), err.Param())

	case "lte":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("Field '%s' must be at most %s characters long",
				err.Field(), err.Param())
		}

		return fmt.Sprintf("Field '%s' must be <= %s", err.Field(), err.Param())

	case "gt":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("Field '%s' must be more than %s characters long",
				err.Field(), err.Param())
		}

		return fmt.Sprintf("Field '%s' must be > %s", err.Field(), err.Param())

	case "lt":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("Field '%s' must be less than %s characters long",
				err.Field(), err.Param())
		}

		return fmt.Sprintf("Field '%s' must be < %s", err.Field(), err.Param())

	case "len":
		return fmt.Sprintf(
			"Field '%s' must be exactly %s characters long", err.Field(), err.Param())

	case "numeric":
		return fmt.Sprintf("Field '%s' must contain only numbers", err.Field())

	case "alpha":
		return fmt.Sprintf("Field '%s' must contain only letters", err.Field())

	case "alphanum":
		return fmt.Sprintf("Field '%s' must contain only letters and numbers", err.Field())

	case "email":
		return fmt.Sprintf("Field '%s' must be in email format", err.Field())

	case "uuid":
		return fmt.Sprintf("Field '%s' must be in UUID format", err.Field())

	case "url":
		return fmt.Sprintf("Field '%s' must be a valid URL", err.Field())

	case "ip":
		return fmt.Sprintf("Field '%s' must be a valid IP address", err.Field())

	// case "oneof": return fmt.Sprintf("Field '%s' must be one of: %s", e.Field(), e.Param())

	case "datetime":
		return fmt.Sprintf("Field '%s' must be in format: %s", err.Field(), err.Param())

	default:
		return err.Error()
	}
}
