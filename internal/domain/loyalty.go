package domain

import "time"

// Пользователь
type User struct {
	ID           int64     `db:"id" json:"id"`
	Login        string    `db:"login" json:"login"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// Баланс пользователя
type Balance struct {
	UserID    int64   `db:"user_id" json:"-"`
	Current   float64 `db:"current" json:"current"`
	Withdrawn float64 `db:"withdrawn" json:"withdrawn"`
}

// Заказ, загруженный пользователем
type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID         int64       `db:"id" json:"-"`
	UserID     int64       `db:"user_id" json:"-"`
	Number     string      `db:"order_number" json:"number"`
	Status     OrderStatus `db:"status" json:"status"`
	Accrual    *float64    `db:"accrual" json:"accrual,omitempty"`
	UploadedAt time.Time   `db:"uploaded_at" json:"uploaded_at"`
}

// Списание средств (вывод)
type Withdrawal struct {
	ID          int64     `db:"id" json:"-"`
	UserID      int64     `db:"user_id" json:"-"`
	Order       string    `db:"order_number" json:"order"`
	Sum         float64   `db:"sum" json:"sum"`
	ProcessedAt time.Time `db:"processed_at" json:"processed_at"`
}
