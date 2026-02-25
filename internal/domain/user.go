package domain

import (
	"time"

	"github.com/google/uuid"
)

// User представляет пользователя системы лояльности
type User struct {
	ID           uuid.UUID `json:"id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserCreateRequest представляет запрос на создание пользователя
type UserCreateRequest struct {
	Login    string `json:"login" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserLoginRequest представляет запрос на аутентификацию
type UserLoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
