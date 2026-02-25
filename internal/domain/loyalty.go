package domain

import (
	"time"

	"github.com/google/uuid"
)

// OperationType представляет тип операции с баллами
type OperationType string

const (
	OperationTypeCredit OperationType = "CREDIT"
	OperationTypeDebit  OperationType = "DEBIT"
)

// LoyaltyOperation представляет операцию с баллами лояльности
type LoyaltyOperation struct {
	ID            uuid.UUID     `json:"id"`
	UserID        uuid.UUID     `json:"user_id"`
	OrderNumber   *string       `json:"order_number,omitempty"`
	OperationType OperationType `json:"operation_type"`
	Amount        float64       `json:"amount"`
	ProcessedAt   time.Time     `json:"processed_at"`
	CreatedAt     time.Time     `json:"created_at"`
}

// BalanceResponse представляет ответ с информацией о балансе
type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

// WithdrawRequest представляет запрос на списание средств
type WithdrawRequest struct {
	Order string  `json:"order" validate:"required"`
	Sum   float64 `json:"sum" validate:"required,gt=0"`
}

// WithdrawalResponse представляет ответ с информацией о списании
type WithdrawalResponse struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
