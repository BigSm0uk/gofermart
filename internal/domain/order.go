package domain

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus представляет статус заказа
type OrderStatus string

const (
	OrderStatusNew         OrderStatus = "NEW"
	OrderStatusProcessing  OrderStatus = "PROCESSING"
	OrderStatusProcessed   OrderStatus = "PROCESSED"
	OrderStatusInvalid     OrderStatus = "INVALID"
)

// Order представляет заказ в системе
type Order struct {
	ID         uuid.UUID   `json:"id"`
	Number     string      `json:"number"`
	UserID     uuid.UUID   `json:"user_id"`
	Status     OrderStatus `json:"status"`
	Accrual    float64     `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// OrderResponse представляет ответ с информацией о заказе
type OrderResponse struct {
	Number     string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    *float64    `json:"accrual,omitempty"`
	UploadedAt string      `json:"uploaded_at"`
}
