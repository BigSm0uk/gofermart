package interfaces

import (
	"context"
	"database/sql"

	"github.com/BigSm0uk/gofermart/internal/domain"
)

type AccrualRepository interface {
	Orders(number int) (*domain.AccrualOrder, error)
	RegisterOrder(order *domain.AccrualOrder) error
	RegisterGood(good *domain.AccrualOrderGood) error
	Ping(ctx context.Context) error
	GetStdDB() *sql.DB
}
