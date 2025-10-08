package handlers

import (
	"github.com/go-playground/validator/v10"
)

// StructValidator реализует интерфейс валидации для Fiber
type StructValidator struct {
	validate *validator.Validate
}

// NewStructValidator создает новый валидатор
func NewStructValidator() *StructValidator {
	return &StructValidator{
		validate: validator.New(),
	}
}

// Validate выполняет валидацию структуры
func (v *StructValidator) Validate(out any) error {
	return v.validate.Struct(out)
}
