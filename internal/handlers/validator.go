package handlers

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type StructValidator struct {
	validate *validator.Validate
}

func NewStructValidator() StructValidator {
	v := validator.New()
	// notblank: строка после TrimSpace должна быть непустой
	_ = v.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		s, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		return strings.TrimSpace(s) != ""
	})
	return StructValidator{validate: v}
}

// Validator needs to implement the Validate method
func (v StructValidator) Validate(out any) error {
	return v.validate.Struct(out)
}
