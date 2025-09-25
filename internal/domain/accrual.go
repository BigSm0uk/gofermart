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
	ID        int64              `db:"id" json:"-"`
	Order     string             `db:"order_number" json:"order"`
	Status    AccrualOrderStatus `db:"status" json:"status"`
	Accrual   *float64           `db:"accrual" json:"accrual,omitempty"`
	CreatedAt time.Time          `db:"created_at" json:"-"`
	Goods     []AccrualOrderGood `json:"goods,omitempty"`
}

// Товары внутри заказа
type AccrualOrderGood struct {
	ID          int64   `db:"id" json:"-"`
	OrderID     int64   `db:"order_id" json:"-"`
	Description string  `db:"description" json:"description"`
	Price       float64 `db:"price" json:"price"`
}

// Правила вознаграждений
type RewardRule struct {
	ID         int64     `db:"id" json:"-"`
	Match      string    `db:"match" json:"match"`
	Reward     float64   `db:"reward" json:"reward"`
	RewardType string    `db:"reward_type" json:"reward_type"` // "%" или "pt"
	CreatedAt  time.Time `db:"created_at" json:"-"`
}
