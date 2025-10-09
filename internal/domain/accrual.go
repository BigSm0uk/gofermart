package domain

import "time"

// Статусы обработки заказов
type AccrualOrderStatus string

const (
	AccrualStatusRegistered AccrualOrderStatus = "REGISTERED"
	AccrualStatusInvalid    AccrualOrderStatus = "INVALID"
	AccrualStatusProcessing AccrualOrderStatus = "PROCESSING"
	AccrualStatusProcessed  AccrualOrderStatus = "PROCESSED"
)

// Заказ для расчёта
type AccrualOrder struct {
	ID        int64              `json:"-"`
	Order     string             `json:"order"`
	Status    AccrualOrderStatus `json:"status"`
	Accrual   *float64           `json:"accrual,omitempty"`
	CreatedAt time.Time          `json:"-"`
	Goods     []AccrualOrderGood `json:"goods,omitempty"`
}

// Товары внутри заказа
type AccrualOrderGood struct {
	ID          int64   `json:"-"`
	OrderID     int64   `json:"-"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// Правила вознаграждений
type RewardRule struct {
	ID         int64     `json:"-"`
	Match      string    `json:"match"`
	Reward     float64   `json:"reward"`
	RewardType string    `json:"reward_type"` // "%" или "pt"
	CreatedAt  time.Time `json:"-"`
}
