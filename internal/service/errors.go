package service

import "errors"

// Ошибки сервисов
var (
	// Ошибки пользователей
	ErrUserRegistrationFailed = errors.New("failed to register user")
	ErrUserLoginFailed        = errors.New("failed to login user")
	
	// Ошибки заказов
	ErrOrderCreationFailed    = errors.New("failed to create order")
	ErrOrderRetrievalFailed   = errors.New("failed to retrieve orders")
	
	// Ошибки баланса
	ErrBalanceRetrievalFailed = errors.New("failed to retrieve balance")
	ErrWithdrawalFailed       = errors.New("failed to process withdrawal")
	ErrWithdrawalHistoryFailed = errors.New("failed to retrieve withdrawal history")
)
