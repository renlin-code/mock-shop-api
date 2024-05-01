package handler

import (
	"fmt"
	"strings"

	"gopkg.in/bluesuncorp/validator.v9"
)

type validationError struct {
	errors []string
}

func (ve validationError) Error() string {
	return fmt.Sprintf("Validations errors: %s", strings.Join(ve.errors, ", "))
}

var v *validator.Validate = validator.New()

func validateInput(data interface{}) error {
	var errors []string
	err := v.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("field %s is %s", strings.ToLower(err.Field()), err.Tag()))
		}
		return &validationError{errors}
	}
	return nil
}
