package middleware

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Validate is a global validator instance
var Validate *validator.Validate

// InitValidator initializes the global validator
func InitValidator() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

// ValidateRequest validates a request struct and returns formatted errors
func ValidateRequest(c echo.Context, request interface{}) error {
	if err := Validate.Struct(request); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  formatValidationError(err),
			Internal: err,
		}
	}
	return nil
}

// formatValidationError formats validator errors into a readable message
func formatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if ok := errors.As(err, &validationErrors); !ok {
		return "Validation failed"
	}

	var messages []string
	for _, e := range validationErrors {
		messages = append(messages, formatFieldError(e))
	}

	if len(messages) == 0 {
		return "Validation failed"
	}

	return messages[0] // Return first error message
}

// formatFieldError formats a single field validation error
func formatFieldError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "Field '" + e.Field() + "' is required"
	case "email":
		return "Field '" + e.Field() + "' must be a valid email address"
	case "min":
		return "Field '" + e.Field() + "' must be at least " + e.Param() + " characters"
	case "max":
		return "Field '" + e.Field() + "' must be at most " + e.Param() + " characters"
	case "url":
		return "Field '" + e.Field() + "' must be a valid URL"
	case "numeric":
		return "Field '" + e.Field() + "' must be a number"
	case "oneof":
		return "Field '" + e.Field() + "' must be one of: " + e.Param()
	default:
		return "Field '" + e.Field() + "' is invalid"
	}
}
