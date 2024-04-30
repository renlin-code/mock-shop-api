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

func ValidateInput(data interface{}) error {
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

type signUpInput struct {
	Name  string `json:"name" validate:"required,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type confirmEmailInput struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,max=30"`
}

type signInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=30"`
}
