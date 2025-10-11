package repo

import "errors"

// Ошибки репозиториев
var (
	// Ошибки базы данных
	ErrDatabaseConnection  = errors.New("database connection failed")
	ErrDatabaseQuery       = errors.New("database query failed")
	ErrDatabaseTransaction = errors.New("database transaction failed")

	// Ошибки пользователей
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserCreationFailed = errors.New("failed to create user")

	// Ошибки заказов
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderCreationFailed = errors.New("failed to create order")

	// Ошибки операций лояльности
	ErrOperationCreationFailed  = errors.New("failed to create loyalty operation")
	ErrBalanceCalculationFailed = errors.New("failed to calculate balance")
)
