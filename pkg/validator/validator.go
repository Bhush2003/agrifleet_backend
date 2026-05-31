package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidationErrors maps field names to error messages.
type ValidationErrors map[string]string

// Validate parses and validates a struct, returning field-level errors.
func Validate(s interface{}) ValidationErrors {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errs := make(ValidationErrors)
	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		errs[field] = fieldMessage(e)
	}
	return errs
}

func fieldMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
	case "uuid4":
		return fmt.Sprintf("%s must be a valid UUID", e.Field())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}
