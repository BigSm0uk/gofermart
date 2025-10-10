package handlers

import "errors"

// Ошибки хендлеров
var (
	// Ошибки валидации запросов
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrMissingOrderNumber = errors.New("order number is required")
	ErrInvalidOrderFormat = errors.New("invalid order number format")
	
	// Ошибки аутентификации
	ErrAuthRequired = errors.New("authorization header required")
	ErrInvalidToken = errors.New("invalid token format")
	ErrExpiredToken = errors.New("invalid or expired token")
)
