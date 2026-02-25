package domain

import "errors"

// Предопределенные ошибки домена
var (
	// Ошибки аутентификации
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUnauthorized       = errors.New("unauthorized")

	// Ошибки заказов
	ErrOrderAlreadyExists      = errors.New("order already exists")
	ErrOrderBelongsToOtherUser = errors.New("order belongs to another user")
	ErrInvalidOrderNumber      = errors.New("invalid order number format")

	// Ошибки баланса
	ErrInsufficientFunds = errors.New("insufficient funds")

	// Ошибки rate limiting
	ErrRateLimited = errors.New("rate limited")

	// Общие ошибки
	ErrNotFound      = errors.New("not found")
	ErrInternalError = errors.New("internal error")
)
