package interfaces

import (
	"context"
	"database/sql"

	"github.com/BigSm0uk/gofermart/internal/domain"
)

type AccrualRepository interface {
	Order(ctx context.Context, number string) (*domain.AccrualOrder, error)
	RegisterOrder(ctx context.Context, order *domain.AccrualOrder) error
	RegisterGood(ctx context.Context, good *domain.AccrualOrderGood) error
	Ping(ctx context.Context) error
	GetStdDB() *sql.DB
}
